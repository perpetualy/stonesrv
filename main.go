package main

import (
	"stonesrv/conf"
	"stonesrv/controllers"
	"stonesrv/database"
	"stonesrv/log"
	"stonesrv/routers"
)

func main() {
	log.Init("stone.log", true)
	//读取配置文件

	conf.Init("")
	database.Init()

	ctrls := []controllers.Controllers{
		&controllers.Register{},
		&controllers.Login{},
		&controllers.Logout{},
		&controllers.UserInfo{},
		&controllers.Update{},
		&controllers.Admin{},
	}
	for _, v := range(ctrls){
		routers.AddController(v)
	}

	routers.RunTLS()
}
