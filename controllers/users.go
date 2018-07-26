package controllers

import (
	"github.com/gin-gonic/gin/json"
	"fmt"
	"net/http"
	"stonesrv/database"
	"stonesrv/log"
	"stonesrv/models"
	"stonesrv/crypto"
	"stonesrv/env"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//Register 注册用户
type Register struct {
	Controllers
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
		context.JSON(400, gin.H{"status": "注册失败，参数错误"})
		return
	}

	key := strings.ToLower(fmt.Sprintf("%s%s%s", req.User, req.Mac, req.Disk0))
	//查找Key是否存在
	//如果已经存在，直接返回
	if database.GetDatabase().IsUserExist(key) {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s] 已经被注册 ", req.User)})
		return
	}

	//验证MAC是否存在
	if database.GetDatabase().IsMACExist(req.Mac) {
		//MAC存在返回
		context.JSON(203, gin.H{"status": "这台机器已经被注册，请联系管理员"})
		return
	}

	//验证DISK0是否存在
	if database.GetDatabase().IsDisk0Exist(req.Disk0) {
		//Disk0存在返回
		context.JSON(203, gin.H{"status": "这台机器已经被注册，请联系管理员"})
		return
	}

	//验证使用时长
	if req.Duration < 0 || req.Duration > 512640 {
		context.JSON(203, gin.H{"status": "非法的使用时长"})
		return
	}

	now := time.Now()
	//expdate := now.AddDate(0, 0, int(req.Duration))
	expdate := now.Add(time.Minute * time.Duration(req.Duration))

	user := req
	user.Key = key
	user.RegDate = now.Format(env.FullDateTimeFormat) //生成注册时间
	user.ExpDate = expdate.Format(env.FullDateTimeFormat) //生成过期时间
	user.Activated = 1

	err = database.GetDatabase().InsertMAC(models.MAC{Key: user.Mac, UserKey: user.Key})
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	err = database.GetDatabase().InsertDisk0(models.Disk0{Key: user.Disk0, UserKey: user.Key})
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	err = database.GetDatabase().InsertUser(user)
	if err != nil {
		context.JSON(203, gin.H{"status": fmt.Sprintf("用户 [%s=] 注册失败", req.User)})
		return
	}

	context.JSON(200, gin.H{"status": "注册成功", "User": user.User})
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
	return p.login
}

//login 登录方法实现
func (p *Login) login(context *gin.Context) {
	log.Info(fmt.Sprintf("login %+v", context))
	// Parse JSON
	loginRequest := models.LoginRequest{}
	if context.Bind(&loginRequest) != nil {
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s%s", loginRequest.User, loginRequest.P1, loginRequest.P2)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil{
		//用户不存在 返回
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，用户不存在"})
		return
	}
	//检验密码
	if strings.Compare(loginRequest.Password, usr.Password) != 0 {
		//密码不正确
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，密码不正确"})
		return
	}
	if usr.Activated != 1{
		//用户已经过期
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，用户已经失效"})
		return			
	}
	loc, lok  := time.LoadLocation("Local")
	regDate, pok := time.ParseInLocation(env.FullDateTimeFormat, usr.RegDate, loc)
	expDate, eok := time.ParseInLocation(env.FullDateTimeFormat, usr.ExpDate, loc)
	if lok!=nil || pok!=nil || eok != nil{
		//获取用户时间出错
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，获取用户时间出错"})
		return
	}
	locExpDate := regDate.Add(time.Minute * time.Duration(usr.Duration))
	now := time.Now()
	if now.After(locExpDate) || now.After(expDate){
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，用户已经过期"})
		return	
	}

	deltaDuration := locExpDate.Sub(now)
	duration := int64(deltaDuration.Minutes())
	
	//生成TOKEN
	token := crypto.GenToken(duration, usr.Salt)

	if strings.Compare(token, "") == 0{
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登录失败，生成TOKEN失败"})
		return
	}

	//客户端需要维护token
	context.JSON(http.StatusOK, gin.H{"status": "登录成功","token":token})
}

//Logout 登出控制器  暂时无用
type Logout struct {
	Controllers
}

//GetGroup 空
func (p *Logout) GetGroup() string {
	return "/auth"
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
	return p.logout
}

//GetFunc 注销方法实现
func (p *Logout) logout(context *gin.Context) {
	log.Info(fmt.Sprintf("logout %+v", context))
	// Parse JSON
	logoutRequest := models.LogoutRequest{}
	if context.Bind(&logoutRequest) != nil {
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登出失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s%s", logoutRequest.User, logoutRequest.P1, logoutRequest.P2)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil{
		//用户不存在 返回
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登出失败，用户不存在"})
		return
	}
	//登出成功 
	token := crypto.GenToken(0, usr.Salt)
	if strings.Compare(token, "") == 0{
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "登出失败"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "登出成功?"})
}

//UserInfo 用户信息控制器
type UserInfo struct {
	Controllers
}

//GetGroup 空
func (p *UserInfo) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 /usr/logout
func (p *UserInfo) GetRelativePath() string {
	return "/usr/info"
}

//GetMethod 方法 POST
func (p *UserInfo) GetMethod() string {
	return "POST"
}

//GetFunc 获取用户信息方法实现
func (p *UserInfo) GetFunc() func(context *gin.Context) {
	return p.getInfo
}

//GetFunc 获取用户信息方法实现
func (p *UserInfo) getInfo(context *gin.Context) {
	log.Info(fmt.Sprintf("getInfo %+v", context))
	// Parse JSON
	userInfoRequest := models.UserInfoRequest{}
	if context.Bind(&userInfoRequest) != nil {
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s%s", userInfoRequest.User, userInfoRequest.P1, userInfoRequest.P2)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil{
		//用户不存在 返回
		context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，用户不存在"})
		return
	}

	//创建回报
	rsp := models.UserInfoResponse{
		User:usr.User,
		FullName:usr.FullName,
		Company:usr.Company,
		Address:usr.Address,
		Email:usr.Email,
		Phone:usr.Phone,
		RegDate:usr.RegDate,
		ExpDate:usr.ExpDate,
	}
	info, err := json.Marshal(rsp)
	if err != nil{
		info = []byte("{}")
	}
	//info := fmt.Sprintf("{\"用户名\":\"%s\",\"全名\":\"%s\",\"公司\":\"%s\",\"地址\":\"%s\",\"电话\":\"%s\",\"邮箱\":\"%s\",\"注册\":\"%s\",\"过期\":\"%s\"}", usr.User, usr.FullName, usr.Company, usr.Address, usr.Phone, usr.Email, usr.RegDate, usr.ExpDate)
	context.JSON(http.StatusOK, gin.H{"status": "获取用户信息成功","info":string(info)})
}
