package controllers

import "github.com/gin-gonic/gin"

//Update 应用版本更新
type Update struct {
	Controllers
}

//GetGroup 空
func (p *Update) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 update
func (p *Update) GetRelativePath() string {
	return "/update"
}

//GetMethod GET
func (p *Update) GetMethod() string {
	return "GET"
}

//GetFunc 更新方法实现
func (p *Update) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		user := context.Params.ByName("name")
		value := user
		ok := true
		if ok {
			context.JSON(200, gin.H{"user": user, "value": value})
		} else {
			context.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	}
}
