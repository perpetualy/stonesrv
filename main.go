package main

import (
	"stonesrv/conf"
	"stonesrv/controllers"
	"stonesrv/database"
	"stonesrv/language"
	"stonesrv/log"
	"stonesrv/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Init("stone.log", true, true)
	//读取配置文件

	conf.Init("")
	language.Init("./language")
	database.Init()

	ctrls := []controllers.Controllers{
		&controllers.Register{},
		&controllers.Login{},
		&controllers.Logout{},
		&controllers.UserInfo{},
		&controllers.Updates{},
		&controllers.Admin{},
	}
	for _, v := range ctrls {
		routers.AddController(v)
	}

	if conf.IsDebugMode() {
		routers.Run()
	} else {
		gin.SetMode(gin.ReleaseMode)
		routers.RunTLS()
	}
}
