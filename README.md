#GMQ 消息模块组合架构

## 自动管理channel 实现配置变动后 模块重启

## 模块重启前保存他的消息队列,重启后在消息原来的消息

## todo list
- 模块停止
- 清空模块队列中的任务
- 模块停止后启动
- 平滑停止,主程序在所有模块的队列中的任务执行完毕后停止
- 讲当前模块中的任务做磁盘持久化


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
 
# mqtt模块配置
```
   mqtt:
      type : inner
      productKey : a1GvAwy4lNt
      deviceSecret: ulsuWZkXGzOQsR5L5AYUQpQxLKvidmUi
      deviceName:  PC1219
```

# web模块功能
