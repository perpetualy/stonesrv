package routers

import (
	"fmt"
	"net/http"
	"path"
	"runtime/debug"
	"stonesrv/accounts"
	"stonesrv/conf"
	"stonesrv/controllers"
	"stonesrv/log"
	"stonesrv/middlewares"
	"sync"

	"github.com/gin-gonic/gin"
)

var rou = newRouters()

func newRouters() *Routers {
	return new(Routers)
}

//Routers 路由
type Routers struct {
	routerEng   *gin.Engine
	controllers sync.Map
}

//AddController 添加控制器
func AddController(controller controllers.Controllers) {
	rou.addController(controller)
}

//Run 运行服务
func Run() {
	rou.run()
}

//RunTLS 运行TLS服务
func RunTLS() {
	rou.runTLS()
}

//regControllers 注册所有控制器
func (p *Routers) regControllers() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("routers.go : regControllers() %v %+v", err, string(debug.Stack())))
		}
	}()
	p.controllers.Range(func(ki, vi interface{}) bool {
		controller := vi.(controllers.Controllers)
		p.regController(p.routerEng, controller)
		return true
	})
}

func (p *Routers) regController(e *gin.Engine, c controllers.Controllers) {
	relativePath := c.GetRelativePath()
	funk := c.GetFunc()
	group := c.GetGroup() //在这里group强制用作认证前置
	method := c.GetMethod()
	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	var irouter gin.IRouter

	switch group {
	case "/":
		irouter = e.Group(group, gin.BasicAuth(accounts.GetAccounts()))
	case "/auth":
		irouter = e.Group(group, middlewares.AuthToken())
	case "":
		irouter = e
	}

	log.Info(fmt.Sprintf("regControllers() %+v %+v %+v %+v", relativePath, &funk, group, method))
	switch method {
	case "POST":
		irouter.POST(relativePath, funk)
	case "GET":
		irouter.GET(relativePath, funk)
	case "DELETE":
		irouter.DELETE(relativePath, funk)
	case "PATCH":
		irouter.PATCH(relativePath, funk)
	case "PUT":
		irouter.PUT(relativePath, funk)
	case "OPTIONS":
		irouter.OPTIONS(relativePath, funk)
	case "HEAD":
		irouter.HEAD(relativePath, funk)
	}
}

func (p *Routers) setStaticFileSystem() {
	//更新文件系统
	relativePath := fmt.Sprintf("%s", conf.GetUpdatesDir())
	dir := path.Base(relativePath)
	p.routerEng.StaticFS(relativePath, http.Dir(dir))
}

func (p *Routers) run() {
	//初始化 Eng
	rou.routerEng = gin.Default()

	//注册控制器
	p.regControllers()

	//设置文件系统
	p.setStaticFileSystem()

	//启动监控
	p.routerEng.Run(fmt.Sprintf("%s:%s", conf.GetServerAddress(), conf.GetServerPort()))
}

func (p *Routers) runTLS() {
	//初始化 Eng
	rou.routerEng = gin.Default()

	//注册控制器
	p.regControllers()

	//设置文件系统
	p.setStaticFileSystem()

	//启动监控
	p.routerEng.RunTLS(fmt.Sprintf("%s:%s", conf.GetServerAddress(), conf.GetServerPort()), conf.GetSSLCrtFile(), conf.GetSSLKeyFile())
}

func (p *Routers) addController(controller controllers.Controllers) {
	p.controllers.Store(controller.GetRelativePath(), controller)
}
