package controllers

import "github.com/gin-gonic/gin"

//应用版本更新
type Update struct {
	Controllers
}

func (p *Update) GetGroup() string{
	return ""
}

func (p *Update) GetRelativePath() string{
	return "update"
}

func (p *Update) GetMethod() string{
	return "GET"
}

func (p *Update) GetFunc() func(context *gin.Context){
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