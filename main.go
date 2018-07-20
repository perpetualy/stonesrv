package main

import (
	"fmt"
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

	////启用获取全部用户线程
	//////////////////////////////////////////////////
	users := database.GetDatabase().GetAllUsers()
	log.Info(fmt.Sprintf("Users:%v", users))
	//////////////////////////////////////////////////

	user := database.GetDatabase().GetUserByName("root")
	log.Info(fmt.Sprintf("User:%v", user))

	routers.Init()
	pingpong := new(controllers.PingPong)
	routers.AddController(pingpong)

	routers.AddController(&controllers.Register{
		Db: database.GetDatabase(),
	})
	routers.AddController(&controllers.Login{})
	routers.AddController(&controllers.Logout{})
	routers.AddController(&controllers.Update{})
	routers.AddController(&controllers.Admin{})

	routers.RunTLS()
}
