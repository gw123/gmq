package server

import (
	"github.com/go-playground/validator"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/webModule/controllers"
	"github.com/gw123/GMQ/modules/webModule/db_models"
	"github.com/gw123/GMQ/modules/webModule/webEvent"
	"github.com/gw123/GMQ/modules/webModule/webMiddlewares"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

type Server struct {
	addr   string
	module interfaces.Module
	echo   *echo.Echo
}

func NewServer(addr string, module interfaces.Module) *Server {
	this := new(Server)
	this.addr = addr
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
			ctx.JSON(http.StatusInternalServerError, map[string]string{"msg": err.Error()})
		}
	}

	publicRoot := this.module.GetConfig().GetStringItem("publicRoot")
	if publicRoot != "" {
		e.Static("/", publicRoot)
	}

	staticHost := this.module.GetConfig().GetStringItem("staticHost")
	version := this.module.GetConfig().GetStringItem("staticVersion")
	viewsRoot := this.module.GetConfig().GetStringItem("viewsRoot")
	e.Renderer = NewTemplateRenderer(viewsRoot, staticHost, version)

	origins := this.module.GetConfig().GetArrayItem("allowOrigins")
	if len(origins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			//authorization,x-csrf-token,x-requested-with
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType,
				echo.HeaderAccept, "x-requested-with", "authorization", "x-csrf-token"},
		}))
	}

	e.Use(webMiddlewares.RequestID())
	timestampFormat := this.module.GetApp().GetConfig().GetString("logger.timestampFormat")
	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05"
	}

	loggerMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_custom}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out},"User-Agent":"${header:User-Agent}",` +
			`"Origin":"${header:Origin}","X-Request-ID":"${header:X-Request-ID}","error":"${error}"}` + "\n",
		//Output: os.Stdout,
		Output: this.module.GetApp(),
		Skipper: func(ctx echo.Context) bool {
			req := ctx.Request()
			return (req.RequestURI == "/" && req.Method == "HEAD") || (req.RequestURI == "/favicon.ico" && req.Method == "GET")
		},
		CustomTimeFormat:timestampFormat,
	})
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:  1 << 10, // 1 KB
	}))

	e.Use(loggerMiddleware)
	//e.Use(webMiddlewares.NewPingMiddleware(this.module.GetApp()))

	indexController := controllers.NewIndexController(this.module)
	taskController := controllers.NewTaskController(this.module)
	commentController := controllers.NewCommentController(this.module)
	resourceController := controllers.NewResourceController(this.module)
	userController := controllers.NewUserController(this.module)

	normalGroup := e.Group("/gapi")
	//normalGroup.Use(webMiddlewares.NewPingMiddleware(this.module.GetApp()))
	normalGroup.GET("/getResource", resourceController.GetResource)
	normalGroup.GET("/getGroup", resourceController.GetGroup)
	normalGroup.GET("/getChapter", resourceController.GetChapter)
	normalGroup.GET("/getCategories", resourceController.GetCategories)
	normalGroup.GET("/getIndexList", resourceController.GetIndexList)

	//GetCategoryCtrl
	normalGroup.GET("/getCategoryCtrl/:category_id/:tag_id", resourceController.GetCategoryCtrl)
	normalGroup.GET("/getCategoryCtrl/:category_id", resourceController.GetCategoryCtrl)

	//登录注册
	normalGroup.POST("/sendMessage", userController.SendMessage)
	normalGroup.POST("/login", userController.Login)
	normalGroup.POST("/register", userController.Register)
	normalGroup.POST("/comments", commentController.CommentList)

	authGroup := e.Group("/gapi")
	authGroup.Use(webMiddlewares.NewAuthMiddleware(this.module))
	authGroup.GET("/getUser", userController.GetUser)
	authGroup.GET("/getUserCollection", userController.GetUserCollection)
	authGroup.POST("/changeCollecton", userController.ChangeUserCollection)

	authGroup.POST("/comment", commentController.Comment)

	authGroup.POST("/addTask", taskController.AddTask)
	authGroup.GET("queryTasksByName", taskController.QueryTasksByName)
	authGroup.POST("/addClientTask", taskController.AddClientTask)

	//client
	clientController := controllers.NewClientController(this.module)
	authGroup.GET("/clientList", clientController.ClientList)
	authGroup.GET("/client/:client_id", clientController.ClientInfo)

	//静态页面
	e.GET("/login", indexController.Login)
	e.GET("/register", indexController.Register)
	e.GET("/", indexController.Index)
	e.GET("/index", indexController.Index)

	e.GET("group/:id", indexController.Group)
	e.GET("resource/:id", indexController.Resource)
	e.GET("chapter/:id", indexController.Chapter)

	e.GET("news", indexController.News)
	e.GET("tagnews/:gid/:gtitle/:tid/:title", indexController.TagNews)
	e.GET("news/:gid/:tid", indexController.News)

	e.GET("home", indexController.Home)
	e.GET("userCollection", indexController.Home)

	e.GET("testpaper/:id", indexController.Testpaper)

	this.module.Info("端口监听在:  %s", this.addr)
	this.echo = e
	return e.Start(this.addr)
}
