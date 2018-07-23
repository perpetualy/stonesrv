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

	routers.Init()
	ctrls := []controllers.Controllers{
		&controllers.Register{
			Db: database.GetDatabase(),
		},
		&controllers.Login{},
		&controllers.Logout{},
		&controllers.Update{},
		&controllers.Admin{},
	}
	for _, v := range(ctrls){
		routers.AddController(v)
	}

	routers.RunTLS()
}
