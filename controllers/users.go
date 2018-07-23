package controllers

import (
	"fmt"
	"stonesrv/accounts"
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
	return func(context *gin.Context) {
		var json struct {
			User     string `json:"user" binding:"required"`
			Password string `json:"password" binding:"required"` //MD5以后的值
			Email    string `json:"email" binding:"required"`
			Address  string `json:"address" binding:"required"`
			FullName string `json:"fullname" binding:"required"`

			Mac   string `json:"mail" binding:"required"`
			Disk0 string `json:"disk0" binding:"required"`
			Salt  int64  `json:"salt" binding:"required"`
		}
		err := context.Bind(&json)
		if err != nil {
			log.Error(fmt.Sprintf("Register JSON error %+v", err))
			context.JSON(8000, gin.H{"status": "Reg failed"})
			return
		}

		//验证MAC是否合法，
		//验证DISK0是否合法
		//如果已经注册过了就不允许注册了。

		now := time.Now()
		nowyear, nowmonth, nowday := now.Date()
		expdate := now.AddDate(0, 1, 0)
		expyear, expmonth, expday := expdate.Date()

		key := strings.ToLower(fmt.Sprintf("%s_%d_%s_%s", json.User, json.Salt, json.Mac, json.Disk0))
		user := models.User{
			Key:      key, //客户端转为小写 + 下划线 + Salt + 下划线 + ENCODE(SALT MAC) + 下划线 + ENCODE(SALT DISK)
			User:     json.User,
			Password: json.Password, // 真实密码 MD5 + ENCODE(SALT)
			Email:    json.Email,
			Address:  json.Address,
			FullName: json.FullName,

			Mac:       json.Mac,   //"1C:1B:0D:A9:3F:79",  //ENCODE(SALT)
			Disk0:     json.Disk0, //"20935UKALWV28LA8923SDF", //ENCODE(SALT)
			Salt:      json.Salt,
			RegDate:   fmt.Sprintf("%d-%d-%d", nowyear, nowmonth, nowday), //生成注册时间
			ExpDate:   fmt.Sprintf("%d-%d-%d", expyear, expmonth, expday), //生成到期时间, 默认试用一个月
			Activated: 1,                                                  //用户默认激活状态
		}
		p.Db.UpsertUser(user)
		accounts.AddAccount(user.User, user.Password)
		context.JSON(200, gin.H{"status": "ok", "User": user.User})
	}
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
		// Parse JSON
		loginRequest := models.LoginRequest{}
		if context.Bind(&loginRequest) != nil {
			context.JSON(203, gin.H{"status": "登录失败，参数不正确"})
			return
		}
		key := fmt.Sprintf("%s%s%s", loginRequest.User, loginRequest.MAC, loginRequest.Disk0)
		usr := database.GetDatabase().GetUserByKey(key)
		//查询用户是否存在
		if usr == nil{
			//用户不存在 返回
			context.JSON(203, gin.H{"status": "登录失败，用户不存在"})
			return
		}
		//检验密码
		if strings.Compare(loginRequest.Password, usr.Password) != 0 {
			//密码不正确
			context.JSON(203, gin.H{"status": "登录失败，密码不正确"})
			return
		}
		if usr.Activated != 1{
			//用户已经过期
			context.JSON(203, gin.H{"status": "登录失败，用户已经过期"})
			return			
		}
		//查找PRIVATE KEY
		//如果PRIVATE KEY存在
		//使用PRIVATE KEY解密 内容 这里需要考虑要解密的内容
		//验证
		//返回SESSIONID和user
		//JSON用PRIVATE KEY加密
		//客户端需要维护SESSION和USER
		context.JSON(200, gin.H{"status": "登录成功","msg":"XXXXXXXYYYYYYYY"})
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
		// Parse JSON
		logoutRequest := models.LogoutRequest{}
		if context.Bind(&logoutRequest) != nil {
			context.JSON(203, gin.H{"status": "登出失败，参数不正确"})
			return
		}
		key := fmt.Sprintf("%s%s%s", logoutRequest.User, logoutRequest.MAC, logoutRequest.Disk0)
		usr := database.GetDatabase().GetUserByKey(key)
		//查询用户是否存在
		if usr == nil{
			//用户不存在 返回
			context.JSON(203, gin.H{"status": "登出失败，用户不存在"})
			return
		}
		//查询用户的SESSIONID
		//删掉服务器对应的SESSIONID
		//登出成功 
		context.JSON(200, gin.H{"status": "登出成功"})
	}
}
