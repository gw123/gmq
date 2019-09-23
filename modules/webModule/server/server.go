package server

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/labstack/echo"
	"github.com/gw123/GMQ/modules/webModule/db_models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/middleware"
	"os"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"fmt"
	"net/http"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"github.com/gw123/GMQ/modules/webModule/webMiddlewares"
)

type Server struct {
	addr   string
	port   int
	module interfaces.Module
	echo   *echo.Echo
}

func NewServer(addr string, port int, module interfaces.Module) *Server {
	this := new(Server)
	this.addr = addr
	this.port = port
	this.module = module
	return this
}

func (this *Server) Start() error {
	e := echo.New()
	e.Validator = &db_models.CustomValidator{Validator: validator.New()}

	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if he, ok := err.(*webEvent.WebError); ok {
			response := he.GetResponse()
			ctx.JSON(http.StatusOK, response)
		} else {
			ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": "内部错误"})
		}
	}

	publicRoot := this.module.GetConfig().GetItem("publicRoot")
	if publicRoot != "" {
		e.Static("/", publicRoot)
	}

	staticFileUrl := this.module.GetConfig().GetItem("staticFileUrl")
	staticFileVersion := this.module.GetConfig().GetItem("staticFileVersion")
	viewsRoot := this.module.GetConfig().GetItem("viewsRoot")
	e.Renderer = NewTemplateRenderer(viewsRoot, staticFileUrl, staticFileVersion)


	loggerMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out},"User-Agent":"${header:User-Agent}",` +
			`"Origin":"${header:Origin}","Content-Type":"${header:Content-Type}","error":"${error}"}` + "\n",
		Output: os.Stdout,
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
	})

	indexController := controllers.NewIndexController(this.module)
	taskController := controllers.NewTaskController(this.module)

	normalGroup := e.Group("")
	normalGroup.Use(webMiddlewares.NewPingMiddleware(this.module.GetApp()))
	normalGroup.GET("/ping", indexController.Ping)
	normalGroup.GET("/status", indexController.Status)
	normalGroup.GET("/testOrder", taskController.CreateTestOrder)

	authGroup := e.Group("")
	authGroup.Use(loggerMiddleware)
	origins := this.module.GetConfig().GetArrayItem("allowOrigins")
	if len(origins) > 0 {
		authGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "x-requested-with", "authorization"},
		}))
	}

	authGroup.GET("/favicon.ico", func(i echo.Context) error {
		return i.String(http.StatusOK, "ico")
	})
	authGroup.GET("/index", indexController.Index)
	authGroup.GET("/message", indexController.Message)

	//登录上报节点有那些服务, 服务版本 , 服务端会下发最新版本下载地址
	authGroup.GET("/login", taskController.Login)
	//下载
	//e.GET("download", taskController.Download)
	authGroup.POST("/addTask", taskController.AddTask)
	authGroup.GET("queryTasksByName", taskController.QueryTasksByName)
	authGroup.POST("/addClientTask", taskController.AddClientTask)

	//client
	clientController := controllers.NewClientController(this.module)
	authGroup.GET("/clientList", clientController.ClientList)
	authGroup.GET("/client/:client_id", clientController.ClientInfo)

	addr := fmt.Sprintf("%s:%d", this.addr, this.port)
	this.module.Info("端口监听在:  %s", addr)
	this.echo = e
	return e.Start(addr)
}
