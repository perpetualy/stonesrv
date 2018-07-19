package controllers

import "github.com/gin-gonic/gin"

//测试端口
type PingPong struct {
	Controllers
}

func (p *PingPong) GetGroup() string{
	return "/"
}

func (p *PingPong) GetRelativePath() string{
	return "ping"
}

func (p *PingPong) GetMethod() string{
	return "GET"
}

func (p *PingPong) GetFunc() func(context *gin.Context){
	return func(context *gin.Context) {
		context.String(200, "pong")
	}
}