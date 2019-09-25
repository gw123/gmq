#GMQ 消息模块组合架构
##  功能特性
  - 自动管理channel 实现配置变动后 模块重启
  - 多模块按需编译,按配置加载
  - 支持多种消息通信方式(mqtt,http,websocket,grpc)
  - 模块之间通过消息之间通信相互独立解耦 
  - 优雅的日志输出方式
  - 多数据库支持,方便开发
  - 支持读取yml中配置环境变量,方便容器化部署
  
## 目录结构
 - bootstarp 引导启动目录, 加载配置, 按需加载需要的module(仅仅编译添加的文件)
 - common 公共引用路径 
 - core 框架核心
 - modules 模块路径
 - resource 资源路径
 - storage 上传文件保存路径


## 如何新加一个model, 下面以debugmodel为例
 - 在modules 添加deubugModel 文件夹
 - 添加DebugModule.go 文件
 ```
import (
	"github.com/gw123/GMQ/modules/base"
	"github.com/gw123/GMQ/core/interfaces"
	"fmt"
)
type DebugModule struct {
	base.BaseModule
}

func NewDebugModule() *DebugModule {
	this := new(DebugModule)
	return this
}

func (this *DebugModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	//订阅 debug 消息主题
	app.Sub("debug", this)
	return nil
}
// 处理消息
func (this *DebugModule) Handle(event interfaces.Event) error {
	return nil
}
// 定时触发方法
func (this *DebugModule) Watch(index int) {

	return
}
```

- 实现模块提供者 DebugModuleProvider.go
```
package debugModule

import "github.com/gw123/GMQ/core/interfaces"

type DebugModuleProvider struct {
	module interfaces.Module
}

func NewDebugModuleProvider() *DebugModuleProvider {
	this := new(DebugModuleProvider)
	return this
}

func (this *DebugModuleProvider) GetModuleName() string {
	return "Debug"
}

func (this *DebugModuleProvider) Register() {
}

func (this *DebugModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewDebugModule()
	return this.module
}

func (this *DebugModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewDebugModule()
	return this.module
}

```

- 在bootstarp/moduleProvider.go 引入debugModuel

```
func LoadModuleProvider(app interfaces.App) {
	app.LoadModuleProvider(debugModule.NewDebugModuleProvider())
	return
}
```



#数据库配置 支持多数据库配置
```
dbpool:
   default: xyt
   db1:
      database: "gateway"
      host: "xytschool.com"
      username: "dbuser"
      password: "dbpwd"
      drive: "mysql,pg,sqllite"
   db2:
      database: "gateway"
      host: "xytschool.com"
      username: "dbuser"
      password: "dbpwd"
      drive: "mysql,pg,sqllite"    
```

#数据库使用
  -  //获取数据库信息
	GetDb(dnname string) (*gorm.DB, error)
  -	//获取默认数据库
	GetDefaultDb() (*gorm.DB, error)
	
	

## 模块使用
```
   moduleName:
      type : inner/exe/dll/so
      enable: 1/ture
```
 moduleName 模块名称
 - type 模块类型目前支持 内部模块(golang),可执行程序, dll/so windows,linux动态库方式
 - 通信方式支持 mqtt,http,grpc,websocket方式
 - enable 配置模块是否要启用
 
## mqtt模块配置 (mqtt 目前使用的是阿里云IOT服务)
```
   mqtt:
      type : inner
      productKey : key
      deviceSecret: secret
      deviceName:  name
```

## web模块功能
```
  web:
      type : inner
      addr: 0.0.0.0
      port: 8080
      staticFileUrl: "http://127.0.0.1:8080"
      staticFileVersion : "1001"
      viewsRoot: "./views"
      publicRoot: "./public"
      sotragePath: "./storage"
      allowOrigins:
          - "http://127.0.0.1:8080"
          - "http://localhost:10086"
          - "http://127.0.0.1:88"
```

## docker服务化支持实现一个评论的服务,具体实现可以参考commentModuel 

### 这里我把编译好多程序命名为gatwway , 下面是gateway的配置文件
```
app:
  logFilterCategories : "EventQueue,Dispath"
  debugLevel:  debug
dbpool:
   default: xyt
   xyt:
      database: "gateway"
      host: "${DB_HOST}"
      username: "${DB_USER}"
      password: "${DB_PWD}"
      drive: "mysql"
modules:
   comment:
     type: inner
     bindAddr: ${COMMENT_ADDR}
```
#### 上面的${var} 的配置是从程序运行的环境变量中读取, 例如 ${DB_HOST}

- mv docker-compose.example.yml docker-compose.yml 修改内容为自己的配置
- 修改 envoy.yaml 配置
- 运行 docker-compose  up -d
 
```
Starting envoy   ... done
Starting gateway ... done
```

## todo list
- 模块停止
- 清空模块队列中的任务
- 平滑停止,主程序在所有模块的队列中的任务执行完毕后停止
- 讲当前模块中的任务做磁盘持久化

## 在线使用案例 httt://m.xyt.xytschool.com 资源组评论和文章评论
