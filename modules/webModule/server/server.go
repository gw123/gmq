package server

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/labstack/echo"
	"github.com/gw123/GMQ/modules/webModule/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/middleware"
	"os"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"fmt"
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
	e.Validator = &models.CustomValidator{Validator: validator.New()}

	loggerMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out},"User-Agent":"${header:User-Agent}",` +
			`"Origin":"${header:Origin}","Content-Type":"${header:Content-Type}","error":"${error}"}` + "\n",
		Output: os.Stdout,
	})
	e.Use(loggerMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	//e.HTTPErrorHandler = func(err error, ctx echo.Context) {
	//	if he, ok := err.(*echo.HTTPError); ok {
	//	}
	//	ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": ""})
	//}

	controller := controllers.NewIndexController(this.module)
	e.GET("/message", controller.Message)
	//e.GET("/", controller.)

	serverController := controllers.NewServerController(this.module)
	e.GET("login", serverController.Login)
	e.GET("download", serverController.Download)
	e.GET("uploadVersion", serverController.UploadVersion)
	e.POST("uploadVersion", serverController.UploadVersion)

	addr := fmt.Sprintf("%s:%d", this.addr, this.port)
	this.module.Info("端口监听在:  %s", addr)
	this.echo = e
	return e.Start(addr)
}
