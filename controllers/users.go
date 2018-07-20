package controllers

import (
	"fmt"
	"stonesrv/database"
	"stonesrv/log"
	"stonesrv/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Register 注册用户
type Register struct {
	Controllers
	Db database.DataBase
}

//GetGroup 空
func (p *Register) GetGroup() string {
	return ""
}

//GetRelativePath 路径/usr/register
func (p *Register) GetRelativePath() string {
	return "/usr/register"
}

//GetMethod 方法 POST
func (p *Register) GetMethod() string {
	return "POST"
}

//GetFunc 注册方法实现
func (p *Register) GetFunc() func(context *gin.Context) {
	return p.register
}

func (p *Register) register(context *gin.Context) {
	req := models.User{}
	err := context.Bind(&req)
	if err != nil {
		log.Error(fmt.Sprintf("Register JSON error %+v", err))
		context.JSON(400, gin.H{"status": "Reg failed"})
		return
	}

	key := strings.ToLower(fmt.Sprintf("%s%s%s", req.User, req.Mac, req.Disk0))
	//查找Key是否存在
	//如果已经存在，直接返回
	if p.Db.IsUserExist(key) {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s] 已经被注册 ", req.User)})
		return
	}

	//验证MAC是否存在
	if p.Db.IsMACExist(req.Mac) {
		//MAC存在返回
		context.JSON(203, gin.H{"status": "这台机器已经被注册，请联系管理员"})
		return
	}

	//验证DISK0是否存在
	if p.Db.IsDisk0Exist(req.Disk0) {
		//Disk0存在返回
		context.JSON(203, gin.H{"status": "这台机器已经被注册，请联系管理员"})
		return
	}

	//验证使用时长
	if req.Duration < 0 || req.Duration > 365 {
		context.JSON(203, gin.H{"status": "非法的使用时长"})
		return
	}

	now := time.Now()
	nowyear, nowmonth, nowday := now.Date()
	expdate := now.AddDate(0, 0, int(req.Duration))
	expyear, expmonth, expday := expdate.Date()

	user := req
	user.Key = key
	user.RegDate = fmt.Sprintf("%d-%d-%d", nowyear, nowmonth, nowday) //生成注册时间
	user.ExpDate = fmt.Sprintf("%d-%d-%d", expyear, expmonth, expday) //生成过期时间
	user.Activated = 1

	err = p.Db.InsertMAC(models.MAC{Key: user.Mac, UserKey: user.Key})
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	err = p.Db.InsertDisk0(models.Disk0{Key: user.Disk0, UserKey: user.Key})
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	err = p.Db.InsertUser(user)
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	context.JSON(200, gin.H{"status": "ok", "User": user.User})
}

//Login 登录控制器
type Login struct {
	Controllers
}

//GetGroup 空
func (p *Login) GetGroup() string {
	return ""
}

//GetRelativePath 路径 /usr/login
func (p *Login) GetRelativePath() string {
	return "/usr/login"
}

//GetMethod 方法 POST
func (p *Login) GetMethod() string {
	return "POST"
}

//GetFunc 登录方法实现
func (p *Login) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		user, _ := context.GetPostForm("user")
		// Parse JSON
		var json struct {
			Value  string `json:"value" binding:"required"`
			MyData string `json:"mydata"`
		}
		//登录后返回TOKEN，用户需要在客户端保存TOKEN
		//token有效期为 10分钟

		if context.Bind(&json) == nil {
			context.JSON(200, gin.H{"status": "ok", "User": user, "Value": json.Value, "MyData": json.MyData})
		}
	}
}

//Logout 登出控制器
type Logout struct {
	Controllers
}

//GetGroup 空
func (p *Logout) GetGroup() string {
	return ""
}

//GetRelativePath 路径 /usr/logout
func (p *Logout) GetRelativePath() string {
	return "/usr/logout"
}

//GetMethod 方法 POST
func (p *Logout) GetMethod() string {
	return "POST"
}

//GetFunc 注销方法实现
func (p *Logout) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		//使用用户名密码可以注销
		context.String(200, "pong")
	}
}
