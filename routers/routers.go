package routers

import (
	"stonesrv/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"stonesrv/accounts"
	"stonesrv/controllers"
	"stonesrv/log"
	"strings"
	"sync"
)

var rou = newRouters()

func newRouters() *Routers{
	return new(Routers)
}

type Routers struct {
	address   string
	routerEng *gin.Engine
	controllers sync.Map
}

//初始化服务地址
//初始化BASIC认证账号
func Init(){
	rou.routerEng = gin.Default()
	rou.address = fmt.Sprintf("%s:%s", conf.GetServerAddress(), conf.GetServerPort())
}

//添加控制器
func AddController(controller controllers.Controllers){
	rou.addController(controller)
}

//运行服务
func Run(){
	rou.run()
}

//运行TLS
func RunTLS(){
	rou.runTLS()
}

//注册所有控制器
func (p *Routers) regControllers(){
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("routers.go : regControllers() %v %+v", err, string(debug.Stack())))
		}
	}()
	p.controllers.Range(func(ki, vi interface{}) bool {
		controller := vi.(controllers.Controllers)
		rp := controller.GetRelativePath()
		f  := controller.GetFunc()
		g  := controller.GetGroup()  //在这里group强制用作认证前置
		m  := controller.GetMethod()
		// Authorized group (uses gin.BasicAuth() middleware)
		// Same than:
		// authorized := r.Group("/")
		// authorized.Use(gin.BasicAuth(gin.Credentials{
		//	  "foo":  "bar",
		//	  "manu": "123",
		//}))
		log.Info(fmt.Sprintf("regControllers() %+v %+v %+v %+v", rp, &f, g, m))
		switch m{
		case "POST":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).POST(rp, f)
			}else{
				p.routerEng.POST(rp, f)
			}
		case "GET":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).GET(rp, f)
			}else{
				p.routerEng.GET(rp, f)
			}
		case "DELETE":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).DELETE(rp, f)
			}else{
				p.routerEng.DELETE(rp, f)
			}
		case "PATCH":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).PATCH(rp, f)
			}else{
				p.routerEng.PATCH(rp, f)
			}
		case "PUT":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).PUT(rp, f)
			}else{
				p.routerEng.PUT(rp, f)
			}
		case "OPTIONS":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).OPTIONS(rp, f)
			}else{
				p.routerEng.OPTIONS(rp, f)
			}
		case "HEAD":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).HEAD(rp, f)
			}else{
				p.routerEng.HEAD(rp, f)
			}
		case "Any":
			if strings.Compare(g, "") != 0{
				p.routerEng.Group(g, gin.BasicAuth(accounts.GetAccounts())).Any(rp, f)
			}else{
				p.routerEng.Any(rp, f)
			}
		}
		return true
	})
}

func (p *Routers) run(){

	//注册控制器
	p.regControllers()

	//启动监控
	p.routerEng.Run(p.address)
}

func (p *Routers) runTLS(){
	//注册控制器
	p.regControllers()

	//启动监控
	p.routerEng.RunTLS(p.address, "server.crt", "server.key")
}

func (p *Routers) addController(controller controllers.Controllers){
	p.controllers.Store(controller.GetRelativePath(), controller)
}



