package main

import (
	"fmt"
	"stonesrv/controllers"
	"stonesrv/log"
	"stonesrv/routers"
	"stonesrv/database"
	"stonesrv/conf"
)

func main() {
	log.Init("stone.log", true)
	//读取配置文件

	////启用获取全部用户线程
	//////////////////////////////////////////////////
	users := database.GetDatabase().GetAllUsers()
	log.Info(fmt.Sprintf("Users:%v", users))
	//////////////////////////////////////////////////

	user := database.GetDatabase().GetUserByName("root")
	log.Info(fmt.Sprintf("User:%v", user))

	routers.Init(conf.GetServerAddress())
	pingpong := new(controllers.PingPong)
	routers.AddController(pingpong)

	routers.AddController(&controllers.Register{
		Db:database.GetDatabase(),
	})
	routers.AddController(&controllers.Login{})
	routers.AddController(&controllers.Logout{})
	routers.AddController(&controllers.Update{})
	routers.AddController(&controllers.Admin{})

	routers.RunTLS()
}
