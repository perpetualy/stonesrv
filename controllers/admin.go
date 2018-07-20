package controllers

import "github.com/gin-gonic/gin"

//应用版本更新
type Admin struct {
	Controllers
}

func (p *Admin) GetGroup() string{
	return "/"
}

func (p *Admin) GetRelativePath() string{
	return "admin"
}

func (p *Admin) GetMethod() string{
	return "POST"
}

//admin login
func (p *Admin) GetFunc() func(context *gin.Context){
	return p.login
}

func (p *Admin) login(context *gin.Context){
	user := context.MustGet(gin.AuthUserKey).(string)
	// Parse JSON
	var json struct {
		Value string `json:"value" binding:"required"`
		MyData string `json:"mydata"`
	}
	//ROOT 用户登录
	//给一个SESSION

	if context.Bind(&json) == nil {
		context.JSON(200, gin.H{"status": "ok", "User":user, "Value":json.Value, "MyData": json.MyData})
	}
}

//禁用用户
type DeActivate struct {
	Controllers
}

func (p *DeActivate) GetGroup() string{
	return "/"
}

func (p *DeActivate) GetRelativePath() string{
	return "admin/deactivate"
}

func (p *DeActivate) GetMethod() string{
	return "POST"
}

func (p *DeActivate) GetFunc() func(context *gin.Context){
	return func(context *gin.Context) {
		//使用用户名和密码
		//并且加上SESSION 来判断是否允许禁用指定用户
		context.String(200, "pong")
	}
}