package controllers

import (
	"github.com/gin-gonic/gin"
)

//Controllers 控制器接口
type Controllers interface {
	//获取分组
	GetGroup() string
	//获取相对路径
	GetRelativePath() string
	//获取方法 POST | GET | PUT | DELETE | Any | 等
	GetMethod() string
	//处理
	GetFunc() func(*gin.Context)
}
