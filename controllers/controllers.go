package controllers

import (
	"fmt"
	"stonesrv/accounts"
	"stonesrv/log"
	"stonesrv/middlewares"

	"github.com/gin-gonic/gin"
)

//Controllers 控制器接口
type Controllers interface {
	//获取分组
	GetGroup() string
	//获取相对路径
	GetRelativePath() string
	//获取方法 POST | GET | PUT | DELETE | Any | 等
	GetMethod() string
	//处理
	GetFunc() func(*gin.Context)
}

func setupRouter(c Controllers) *gin.Engine {
	r := gin.Default()
	regController(r, c)
	return r
}

func regController(e *gin.Engine, c Controllers) {
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
	case "/pack":
		irouter = e.Group(group, middlewares.PackToken())
	case "/space":
		irouter = e.Group(group, middlewares.SpaceToken())
	case "/table":
		irouter = e.Group(group, middlewares.TableToken())
	case "/spacetable":
		irouter = e.Group(group, middlewares.SpaceTableToken())
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
