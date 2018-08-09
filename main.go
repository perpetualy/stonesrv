package main

import (
	"stonesrv/conf"
	"stonesrv/controllers"
	"stonesrv/database"
	"stonesrv/language"
	"stonesrv/log"
	"stonesrv/routers"
)

func main() {
	log.Init("stone.log", true)
	//读取配置文件

	conf.Init("")
	language.Init()
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

	routers.RunTLS()
}
