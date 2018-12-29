# erp-client-s 消息总线

## App 订阅事件 configChange , stopModule ,startModule

## 模块常用方法
-  获取模块配置 this := NewTestModule() this是模块的一个实例
```
      dllpath  := this.Config.GetItem("path")
      isInner := this.Config.GetItem("inner")

```

-   获取全局配置
```
      version := this.Config.GetGlobalItem("version")
      basePath := this.Config.GetGlobalItem("basePath")

```

-   订阅事件
```
      this.Sub("vioce")
```

-   取消订阅
```
      this.UnSub("vioce")
```

## 模块配置
- subs: 模块订阅事件 
- inner: 是否内置模块
- path: 外部模块加载路径
- enable: 是否加载模块
- 其他配置

## 模块接口
```
//模块启动函数,第一个参数是配置json字符串,第二个参数是pub函数指针
start([]byte ,func PubCallback(event))

//模块停止函数
stop()

//事件处理函数
handel(event)

//
```

- handel(event) ,event 参数举例
```
    {
            "EventName":"voice",
            "MsgId": "123456000010",
            "Payload": {
                "timestamp" : 1291229212,
                "expired" : 1000,
                "msgId" : "123456",
                "deviceNames" : [],
                "event" : "voice",
                "data" : "123",
            }
    }
```

- pubCallback(event),event 参数举例
 ```
   {
        "EventName":"replay",
        "MsgId": "123456",
        "Payload": {
                        "timestamp" : 1291229212,
                        "expired" : 1000,
                        "msgId" : "123456",
                        "deviceNames" : [],
                        "event" : "voice",
                        "data" : "播放成功:123",
                   }
   }
 ```

### 事件
- sendMqttMsg
  ```
    {
         "EventName":"sendMqttMsg",
         "MsgId": "123456",
         "Payload": {
                         "timestamp": 1231231
                         "msgId" : "123456",
                         "event" : "log",
                         "data" : "日志内容",
                    }
    }
  ```

- replay   *replay事件实例可以使用NewReplyEvent()生成 ,注意NewReplyEvent的Payload的内容必须是LhMsg的json字符串否则服务端无法解析.
  ```
    {
         "EventName":"sendMqttMsg",
         "MsgId": "123456",
         "Payload": "{
                      'timestamp' : 1291229212,
                      'expired' : 1000,
                      'msgId' : "123456",
                      'deviceNames' : [],
                      'event' : "voice",
                      'data' : "{}",
                     }"
    }
  ```

###  mqtt 模块上报消息名称

-                   本机事件名    服务端事件名
- 上报消息处理结果       => 和下发相同
- 上报日志消息     sendMqttMsg =>
- 上报版本                     =>

|事件名|服务端名称|事件说明|
| --------   | -----:   | :----: |
| reply   | 和下发事件名相同  | 上报消息处理结果 |
| sendMqttMsg | log   | 上报日志消息 |
| version  | version  | 上报版本 |
| voice  | voice  | 语言播报 |
| update  | update  | 升级 |








