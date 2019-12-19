package webSocketModule

import (
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/modules/grpcModule/grpcModel"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net"
	"errors"
	"time"
)

type Gserver struct {
	stream     grpcModel.Conn_RWStreamServer
	isStreamOk bool
	module     *GrpcModule
	ctlChan    chan int
}

func NewGserver(module *GrpcModule, msgChan chan interfaces.Msg) *Gserver {
	if module == nil {
		return nil
	}
	this := new(Gserver)
	this.module = module
	this.ctlChan = make(chan int, 1)
	this.isStreamOk = false
	this.stream = nil
	return this
}

func (this *Gserver) RWStream(stream grpcModel.Conn_RWStreamServer) error {
	this.isStreamOk = true
	this.stream = stream
	ctx := stream.Context()
	cancelCtx, cancelFun := context.WithCancel(ctx)
	go func() {
		defer func() {
			this.ctlChan <- 0
			cancelFun()
		}()
		for {
			select {
			case <-cancelCtx.Done():
				return
			default:
				break
			}

			msg, err := stream.Recv()
			if err == io.EOF {
				this.module.Warning("stream recv :" + "客户端发送数据结束")
				return
			}
			if err != nil {
				this.module.Warning("stream recv :" + err.Error())
				return
			}
			event := gmsg.NewEvent(msg.EventName, []byte(msg.Payload))
			this.module.Push(event)
			//fmt.Println(msg.ModuleName, msg.EventName, msg.MsgId, msg.Payload)
			time.Sleep(time.Millisecond)
		}
	}()
	return nil
}

func (this *Gserver) Server() {

}

func (this *Gserver) Push(msg2 interfaces.Msg) error {
	if this.stream == nil {
		return errors.New("client not connect")
	}
	if !this.isStreamOk {
		return errors.New("stream is no ok")
	}
	msg := &grpcModel.Msg{
		ModuleName: msg2.GetSourceModule(),
		MsgId:      msg2.GetMsgId(),
		Payload:    string(msg2.GetPayload()),
		EventName:  msg2.GetEventName(),
	}
	this.stream.Send(msg)
	return nil
}

type GrpcModule struct {
	base.BaseModule
	eventNames    []string
	Port          string
	authToken     string
	clientName    string
	Conn          grpcModel.ConnServer
	Server        *grpc.Server
	gServer       *Gserver
	isServerStart bool
}

func NewGrpcModule() *GrpcModule {
	this := new(GrpcModule)
	return this
}

func (this *GrpcModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	this.Port = config.GetStringItem("port")
	this.gServer = &Gserver{}
	this.Conn = this.gServer
	this.gServer.module = this
	go this.Start()
	return nil
}

func (this *GrpcModule) UnInit() error {
	this.BaseModule.UnInit()
	this.Server.Stop()
	this.isServerStart = false
	return nil
}

func (this *GrpcModule) GetStatus() uint64 {
	return 1
}

func (this *GrpcModule) Handle(event interfaces.Msg) error {
	err := this.gServer.Push(event)
	return err
}

func (this *GrpcModule) Watch(index int) {

}

func (this *GrpcModule) InitGrpc() error {
	lis, err := net.Listen("tcp", this.Port)
	if err != nil {
		this.Error("监听失败" + err.Error())
		return err
	}
	s := grpc.NewServer()
	grpcModel.RegisterConnServer(s, this.Conn)
	s.Serve(lis)
	this.Server = s
	this.isServerStart = true
	return nil
}
