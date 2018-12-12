package main

import (
	"fmt"
	"stonesrv/conf"
	"stonesrv/controllers"
	"stonesrv/database"
	"stonesrv/language"
	"stonesrv/log"
	"stonesrv/routers"

	"github.com/gin-gonic/gin"
)

var _VERSION_ = "0.5.0"

func main() {
	log.Init("stone.log", true, true)
	//读取配置文件

	//show version
	log.Info(fmt.Sprintf("Stone Version %v", _VERSION_))

	conf.Init("")
	language.Init("./language")
	database.Init()

	ctrls := []controllers.Controllers{
		&controllers.Register{},
		&controllers.Login{},
		&controllers.Logout{},
		&controllers.UserInfo{},
		&controllers.LoginWeChat{},
		&controllers.LogoutWeChat{},
		&controllers.UserInfoWeChat{},
		&controllers.Updates{},
		&controllers.Admin{},
		&controllers.InsertPack{},
		&controllers.GetPackToken{},
		&controllers.GetPack{},
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
