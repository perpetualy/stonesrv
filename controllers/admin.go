package controllers

import "github.com/gin-gonic/gin"

//Admin 管理员
type Admin struct {
	Controllers
}

//GetGroup 分组 /
func (p *Admin) GetGroup() string {
	return "/"
}

//GetRelativePath 路径 admin
func (p *Admin) GetRelativePath() string {
	return "admin"
}

//GetMethod 方法 POST
func (p *Admin) GetMethod() string {
	return "POST"
}

//GetFunc admin登录
func (p *Admin) GetFunc() func(context *gin.Context) {
	return p.login
}

//loing admin登录实现
func (p *Admin) login(context *gin.Context) {
	user := context.MustGet(gin.AuthUserKey).(string)
	// Parse JSON
	var json struct {
		Value  string `json:"value" binding:"required"`
		MyData string `json:"mydata"`
	}
	//ROOT 用户登录
	//给一个SESSION

	if context.Bind(&json) == nil {
		context.JSON(200, gin.H{"status": "ok", "User": user, "Value": json.Value, "MyData": json.MyData})
	}
}

//DeActivate 禁用用户
type DeActivate struct {
	Controllers
}

//GetGroup 分组 /
func (p *DeActivate) GetGroup() string {
	return "/"
}

//GetRelativePath 路径 admin/deactivate
func (p *DeActivate) GetRelativePath() string {
	return "admin/deactivate"
}

//GetMethod 方法 POST
func (p *DeActivate) GetMethod() string {
	return "POST"
}

//GetFunc 禁用用户实现
func (p *DeActivate) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		//使用用户名和密码
		//并且加上SESSION 来判断是否允许禁用指定用户
		context.String(200, "pong")
	}
}
