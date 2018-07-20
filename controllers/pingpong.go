package controllers

import "github.com/gin-gonic/gin"

//PingPong 测试
type PingPong struct {
	Controllers
}

//GetGroup 分组 /
func (p *PingPong) GetGroup() string {
	return "/"
}

//GetRelativePath 路径 ping
func (p *PingPong) GetRelativePath() string {
	return "ping"
}

//GetMethod 方法 GET
func (p *PingPong) GetMethod() string {
	return "GET"
}

//GetFunc 测试实现
func (p *PingPong) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.String(200, "pong")
	}
}
