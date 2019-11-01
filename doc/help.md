# grpc 配置代码分析
grpc server 所有配置项目
``` 
type serverOptions struct {
	creds                 credentials.TransportCredentials
	codec                 baseCodec
	cp                    Compressor
	dc                    Decompressor
	unaryInt              UnaryServerInterceptor
	streamInt             StreamServerInterceptor
	inTapHandle           tap.ServerInHandle
	statsHandler          stats.Handler
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	unknownStreamDesc     *StreamDesc
	keepaliveParams       keepalive.ServerParameters
	keepalivePolicy       keepalive.EnforcementPolicy
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	connectionTimeout     time.Duration
	maxHeaderListSize     *uint32
}

//ServerOption 默认 server 配置
var defaultServerOptions = serverOptions{
	maxReceiveMessageSize: defaultServerMaxReceiveMessageSize,
	maxSendMessageSize:    defaultServerMaxSendMessageSize,
	connectionTimeout:     120 * time.Second,
	writeBufferSize:       defaultWriteBufSize,
	readBufferSize:        defaultReadBufSize,
}
```


## ServerOption 接口
```go
// ServerOption 接口 ,可以理解为一个包装类,用来包装 serverOptions
type ServerOption interface {
	apply(*serverOptions)
}
```

## funcServerOption 生成一个包装函数
```go
// funcServerOption 是一个包装函数 ,它实现了ServerOption可以修改serverOptions内容.
type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}
```

## funcServerOption包装函数实例

### 修改grpc server发送缓存区大小
```go
func WriteBufferSize(s int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.writeBufferSize = s
	})
}
```

### 修改保活选项
```go
func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	if kp.Time > 0 && kp.Time < time.Second {
		grpclog.Warning("Adjusting keepalive ping interval to minimum period of 1s")
		kp.Time = time.Second
	}

	return newFuncServerOption(func(o *serverOptions) {
		o.keepaliveParams = kp
	})
}
```

### 设置自己的解码器
```go
func CustomCodec(codec Codec) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.codec = codec
	})
}
```

### 设置最大的接受消息大小
```go
func MaxRecvMsgSize(m int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.maxReceiveMessageSize = m
	})
}

```

### 拦截器(中间件),只能配置一个拦截器, 多个拦截器可以在这个拦截器里面实现.
```go
func UnaryInterceptor(i UnaryServerInterceptor) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		if o.unaryInt != nil {
			panic("The unary server interceptor was already set and may not be reset.")
		}
		o.unaryInt = i
	})
}
```


####详细说明grpc拦截器
```go
type UnaryClientInterceptor 
            func(ctx context.Context, 
			method string, 
			req, reply interface{},
			 cc *ClientConn, 
			invoker UnaryInvoker, 
			opts ...CallOption) error
```

## 配置日志中间件
```go
grpcServer := grpc.NewServer(
    grpc_middleware.WithUnaryServerChain(
        grpc_logrus.UnaryServerInterceptor(entry),
        grpc_ctxtags.UnaryServerInterceptor(
            grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
            grpc_ctxtags.WithFieldExtractor(func(fullMethod string, req interface{}) map[string]interface{} {
                return map[string]interface{}{ "requestData": req}
            }),
        ),
))
```

### 下面分析 ctx_logrus 具体代码实现
ctx_logrus 暴露了三个比较实用的方法 AddFields, Extract ,ToContext,下面看看他们的具体实现
```go
package ctxlogrus

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)
type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger *logrus.Entry
	fields logrus.Fields
}

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx context.Context, entry *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger: entry,
		fields: logrus.Fields{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

//添加日志字段到日志中间件(ctx_logrus)
func AddFields(ctx context.Context, fields logrus.Fields) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}
//导出ctx_logrus日志库和grpc_ctxtags中间件在中记录的日志内容
// Extract takes the call-scoped logrus.Entry from ctx_logrus middleware.
// If the ctx_logrus middleware wasn't used, a no-op `logrus.Entry` is returned. This makes it safe to
// use regardless.
func Extract(ctx context.Context) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(nullLogger)
	}

	fields := logrus.Fields{}

	// Add grpc_ctxtags tags metadata until now.
	tags := grpc_ctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		fields[k] = v
	}

	// Add logrus fields added until now.
	for k, v := range l.fields {
		fields[k] = v
	}

	return l.logger.WithFields(fields)
}
```

### grpc_logrus.UnaryServerInterceptor() 如何将 ctxlogrus context 导入到 context
```go
package grpc_logrus

import (
	"path"
	"time"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags/logrus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	// SystemField is used in every log statement made through grpc_logrus. Can be overwritten before any initialization code.
	SystemField = "system"

	// KindField describes the log gield used to incicate whether this is a server or a client log statment.
	KindField = "span.kind"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(entry *logrus.Entry, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, entry, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		if !o.shouldLog(info.FullMethod, err) {
			return resp, err
		}
		code := o.codeFunc(err)
		level := o.levelFunc(code)
		durField, durVal := o.durationFunc(time.Since(startTime))
		fields := logrus.Fields{
			"grpc.code": code.String(),
			durField:    durVal,
		}
		if err != nil {
			fields[logrus.ErrorKey] = err
		}

		levelLogf(
			ctx_logrus.Extract(newCtx).WithFields(fields), // re-extract logger from newCtx, as it may have extra fields that changed in the holder.
			level,
			"finished unary call with code "+code.String())

		return resp, err
	}
}

func levelLogf(entry *logrus.Entry, level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.DebugLevel:
		entry.Debugf(format, args...)
	case logrus.InfoLevel:
		entry.Infof(format, args...)
	case logrus.WarnLevel:
		entry.Warningf(format, args...)
	case logrus.ErrorLevel:
		entry.Errorf(format, args...)
	case logrus.FatalLevel:
		entry.Fatalf(format, args...)
	case logrus.PanicLevel:
		entry.Panicf(format, args...)
	}
}
// 下面是关键地方,调用上面ctxlogrus.ToContext 将用户配置的entry *logrus.Entry
//导入到 ctxlogrus context中
func newLoggerForCall(ctx context.Context, entry *logrus.Entry, 
        fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	callLog := entry.WithFields(
		logrus.Fields{
			SystemField:       "grpc",
			KindField:         "server",
			"grpc.service":    service,
			"grpc.method":     method,
			"grpc.start_time": start.Format(time.RFC3339),
		})

	if d, ok := ctx.Deadline(); ok {
		callLog = callLog.WithFields(
			logrus.Fields{
				"grpc.request.deadline": d.Format(time.RFC3339),
			})
	}

	callLog = callLog.WithFields(ctx_logrus.Extract(ctx).Data)
	return ctxlogrus.ToContext(ctx, callLog)
}
```


### 代码实例在request中间中将requestid 记录到日志中去
```go
//request中间件 将request_id 记录到请求日志中日志中
func request(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	var requestId string
	if val, ok := md[constants.RequestId]; ok {
		requestId = val[0]
	} else {
		return ctx, status.Errorf(codes.Unauthenticated, "no metadata %s", constants.RequestId)
	}

	ctx = context.WithValue(ctx, constants.RequestId, requestId)
    //将请求的request_id 记录到日志中间件ctxlogrus中
    //这样在本次请求的其他的日志都可以获取到RequestId,方便日志跟踪
	ctxlogrus.AddFields(ctx, logrus.Fields{
		constants.RequestId: requestId,
	})
	return ctx, nil
}

func Request() grpc.UnaryServerInterceptor {
	interceptor := func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, err = request(ctx)
		if err != nil {
			return
		}
		return handler(ctx, req)
	}
	return interceptor
}

```

```go
// 在请求的其他地方可以使用ctxlogrus.Extract(ctx)导出*logrus.Entry来记录日志
// ctx保存了request_id 下面的日志会输出request_id
ctxlogrus.Extract(ctx).Errorf("数据库错误: %s", err.Error())
```
