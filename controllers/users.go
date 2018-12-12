package controllers

import (
	"fmt"
	"runtime/debug"
	"stonesrv/crypto"
	"stonesrv/database"
	"stonesrv/env"
	"stonesrv/log"
	"stonesrv/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin/json"

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
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : register() %v %+v", err, string(debug.Stack())))
		}
	}()
	req := models.User{}
	err := context.Bind(&req)
	if err != nil {
		log.Error(fmt.Sprintf("Register JSON error %+v", err))
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(env.RegFailedParamsErrors, gin.H{"status": language.GetText(env.RegFailedParamsErrors)})
		return
	}

	key := strings.ToLower(fmt.Sprintf("%s%s%s", req.User, req.Mac, req.Disk0))
	//查找Key是否存在
	//如果已经存在，直接返回
	if database.GetDatabase().IsUserExist(key) {
		env.GenJSONResponse(context, env.RegFailedUserAlreadyRegistered, req.User)
		//context.JSON(env.RegFailedUserAlreadyRegistered, gin.H{"status": fmt.Sprintf(language.GetText(env.RegFailedUserAlreadyRegistered), req.User)})
		return
	}

	// //验证MAC是否存在
	// if database.GetDatabase().IsMACExist(req.Mac) {
	// 	//MAC存在返回
	// 	env.GenJSONResponse(context, env.RegFailedPCAlreadyRegistered, nil)
	// 	//context.JSON(env.RegFailedPCAlreadyRegistered, gin.H{"status": language.GetText(env.RegFailedPCAlreadyRegistered)})
	// 	return
	// }

	//验证DISK0是否存在
	if database.GetDatabase().IsDisk0Exist(req.Disk0) {
		//Disk0存在返回
		env.GenJSONResponse(context, env.RegFailedPCAlreadyRegistered, nil)
		//context.JSON(env.RegFailedPCAlreadyRegistered, gin.H{"status": language.GetText(env.RegFailedPCAlreadyRegistered)})
		return
	}

	//验证使用时长 不得是负数 最长年限不能超过10年
	if req.Duration < 0 || req.Duration > 5184000 {
		env.GenJSONResponse(context, env.RegFailedInvalidDuration, nil)
		//context.JSON(env.RegFailedInvalidDuration, gin.H{"status": language.GetText(env.RegFailedInvalidDuration)})
		return
	}

	now := time.Now()
	//expdate := now.AddDate(0, 0, int(req.Duration))
	expdate := now.Add(time.Minute * time.Duration(req.Duration))

	user := req
	user.Key = key
	user.RegDate = now.Format(env.FullDateTimeFormat)     //生成注册时间
	user.ExpDate = expdate.Format(env.FullDateTimeFormat) //生成过期时间
	user.Activated = 1

	unixoffest := time.Now().Unix()
	err = database.GetDatabase().InsertMAC(models.MAC{Key: fmt.Sprintf("%s_%v", user.Mac, unixoffest), UserKey: user.Key})
	if err != nil {
		//返回JSON
		env.GenJSONResponse(context, env.RegFailed, req.User)
		//context.JSON(env.RegFailed, gin.H{"status": fmt.Sprintf(language.GetText(env.RegFailed), req.User)})
		return
	}

	err = database.GetDatabase().InsertDisk0(models.Disk0{Key: user.Disk0, UserKey: user.Key})
	if err != nil {
		//删除多余MAC
		database.GetDatabase().RemoveMAC(models.MAC{Key: user.Mac})
		//返回JSON
		env.GenJSONResponse(context, env.RegFailed, req.User)
		//context.JSON(env.RegFailed, gin.H{"status": fmt.Sprintf(language.GetText(env.RegFailed), req.User)})
		return
	}

	err = database.GetDatabase().InsertUser(user)
	if err != nil {
		//删除多余MAC
		database.GetDatabase().RemoveMAC(models.MAC{Key: user.Mac})
		//删除多余Disk0
		database.GetDatabase().RemoveDisk0(models.Disk0{Key: user.Disk0})
		//返回JSON
		env.GenJSONResponse(context, env.RegFailed, req.User)
		//context.JSON(env.RegFailed, gin.H{"status": fmt.Sprintf(language.GetText(env.RegFailed), req.User)})
		return
	}
	//返回JSON
	env.GenJSONResponse(context, env.RegSuccess, req.User)
	//context.JSON(env.RegSuccess, gin.H{"status": language.GetText(env.RegSuccess)})
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
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : login() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("login %+v", context))
	// Parse JSON
	loginRequest := models.LoginRequest{}
	if context.Bind(&loginRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(env.LoginFailedParamsErrors, gin.H{"status": language.GetText(env.LoginFailedParamsErrors)})
		return
	}
	key := fmt.Sprintf("%s%s%s", loginRequest.User, loginRequest.P1, loginRequest.P2)
	log.Info(key)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.LoginFailedUserDoNotExists, nil)
		//context.JSON(env.LoginFailedUserDoNotExists, gin.H{"status": language.GetText(env.LoginFailedUserDoNotExists)})
		return
	}
	//检验密码
	if strings.Compare(loginRequest.Password, usr.Password) != 0 {
		//密码不正确
		env.GenJSONResponse(context, env.LoginFailedPasswordIncorrect, nil)
		//这里记录密码错误登录失败
		database.GetDatabase().RecordUserPasswordFailed(key)
		//context.JSON(env.LoginFailedPasswordIncorrect, gin.H{"status": language.GetText(env.LoginFailedPasswordIncorrect)})
		return
	}
	if usr.Activated != 1 {
		//用户已经失效
		env.GenJSONResponse(context, env.LoginFailedUserInactivated, nil)
		//这里记录用户失效
		database.GetDatabase().RecordUserInActivated(key)
		//context.JSON(env.LoginFailedUserInactivated, gin.H{"status": language.GetText(env.LoginFailedUserInactivated)})
		return
	}
	loc, lok := time.LoadLocation("Local")
	regDate, pok := time.ParseInLocation(env.FullDateTimeFormat, usr.RegDate, loc)
	expDate, eok := time.ParseInLocation(env.FullDateTimeFormat, usr.ExpDate, loc)
	if lok != nil || pok != nil || eok != nil {
		//获取用户时间出错
		env.GenJSONResponse(context, env.LoginFailedGetDateFailed, nil)
		//context.JSON(env.LoginFailedGetDateFailed, gin.H{"status": language.GetText(env.LoginFailedGetDateFailed)})
		return
	}
	locExpDate := regDate.Add(time.Minute * time.Duration(usr.Duration))
	now := time.Now()
	if now.After(locExpDate) || now.After(expDate) {
		//用户已经过期
		env.GenJSONResponse(context, env.LoginFailedUserExpired, nil)
		//这里记录用户过期
		database.GetDatabase().RecordUserExpired(key)
		//context.JSON(env.LoginFailedUserExpired, gin.H{"status": language.GetText(env.LoginFailedUserExpired)})
		return
	}

	deltaDuration := locExpDate.Sub(now)
	duration := int64(deltaDuration.Minutes())

	//生成TOKEN
	token := crypto.GenToken(duration, usr.Salt)

	if strings.Compare(token, "") == 0 {
		env.GenJSONResponse(context, env.LoginFailedGenTokenFailed, nil)
		//context.JSON(env.LoginFailedGenTokenFailed, gin.H{"status": language.GetText(env.LoginFailedGenTokenFailed)})
		return
	}

	//这里记录登录成功次数
	database.GetDatabase().RecordUserLoginSuccess(key)
	//这里记录登录IP
	//database.GetDatabase().RecordUserLoginIP(key, "IP")
	//客户端需要维护token
	env.GenJSONResponseWithMsg(context, env.LoginSuccess, token)
	//context.JSON(env.LoginSuccess, gin.H{"status": token})
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
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : logout() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("logout %+v", context))
	// Parse JSON
	logoutRequest := models.LogoutRequest{}
	if context.Bind(&logoutRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(env.LogoutFailedParamsErrors, gin.H{"status": "登出失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s%s", logoutRequest.User, logoutRequest.P1, logoutRequest.P2)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.LogoutFailedUserDoNotExists, nil)
		//context.JSON(env.LogoutFailedUserDoNotExists, gin.H{"status": "登出失败，用户不存在"})
		return
	}
	//登出成功
	token := crypto.GenToken(0, usr.Salt)
	if strings.Compare(token, "") == 0 {
		env.GenJSONResponse(context, env.LogoutFailed, nil)
		//context.JSON(env.LogoutFailed, gin.H{"status": "登出失败"})
		return
	}
	//这里记录登出时间
	database.GetDatabase().RecordUserLogoutSuccess(key)
	env.GenJSONResponse(context, env.LogoutSuccess, nil)
	//context.JSON(env.LogoutSuccess, gin.H{"status": "登出成功?"})
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
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : getInfo() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("getInfo %+v", context))
	// Parse JSON
	userInfoRequest := models.UserInfoRequest{}
	if context.Bind(&userInfoRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s%s", userInfoRequest.User, userInfoRequest.P1, userInfoRequest.P2)
	log.Info(key)
	usr := database.GetDatabase().GetUserByKey(key)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.GetUserInfoFailedUserDoNotExists, nil)
		//context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，用户不存在"})
		return
	}

	//创建回报
	rsp := models.UserInfoResponse{
		User:      usr.User,
		FullName:  usr.FullName,
		Company:   usr.Company,
		Address:   usr.Address,
		Email:     usr.Email,
		Phone:     usr.Phone,
		Space:     usr.Space,
		Tables:    usr.Tables,
		Functions: usr.Functions,
		RegDate:   usr.RegDate,
		ExpDate:   usr.ExpDate,
	}
	info, err := json.Marshal(rsp)
	if err != nil {
		info = []byte("{}")
		env.GenJSONResponse(context, env.GetUserInfoFailed, nil)
		return
	}
	env.GenJSONResponseWithMsg(context, env.GetUserInfoSuccess, string(info))
	//info := fmt.Sprintf("{\"用户名\":\"%s\",\"全名\":\"%s\",\"公司\":\"%s\",\"地址\":\"%s\",\"电话\":\"%s\",\"邮箱\":\"%s\",\"注册\":\"%s\",\"过期\":\"%s\"}", usr.User, usr.FullName, usr.Company, usr.Address, usr.Phone, usr.Email, usr.RegDate, usr.ExpDate)
	//context.JSON(http.StatusOK, gin.H{"status": "获取用户信息成功", "info": string(info)})
}

//以下供微信使用
//LoginWeChat 登录控制器
type LoginWeChat struct {
	Controllers
}

//GetGroup 空
func (p *LoginWeChat) GetGroup() string {
	return ""
}

//GetRelativePath 路径 /usr/login
func (p *LoginWeChat) GetRelativePath() string {
	return "/usr/loginwechat"
}

//GetMethod 方法 POST
func (p *LoginWeChat) GetMethod() string {
	return "POST"
}

//GetFunc 登录方法实现
func (p *LoginWeChat) GetFunc() func(context *gin.Context) {
	return p.loginwechat
}

//loginwechat 登录方法实现
func (p *LoginWeChat) loginwechat(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : loginwechat() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("login %+v", context))
	// Parse JSON
	loginRequest := models.LoginRequest{}
	if context.Bind(&loginRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(env.LoginFailedParamsErrors, gin.H{"status": language.GetText(env.LoginFailedParamsErrors)})
		return
	}
	key := fmt.Sprintf("%s%s", loginRequest.User, loginRequest.P1)
	log.Info(key)
	usr := database.GetDatabase().GetUserByName(loginRequest.User)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.LoginFailedUserDoNotExists, nil)
		//context.JSON(env.LoginFailedUserDoNotExists, gin.H{"status": language.GetText(env.LoginFailedUserDoNotExists)})
		return
	}
	//检验密码
	if strings.Compare(loginRequest.Password, usr.Password) != 0 {
		//密码不正确
		env.GenJSONResponse(context, env.LoginFailedPasswordIncorrect, nil)
		//这里记录密码错误登录失败
		database.GetDatabase().RecordUserPasswordFailed(key)
		//context.JSON(env.LoginFailedPasswordIncorrect, gin.H{"status": language.GetText(env.LoginFailedPasswordIncorrect)})
		return
	}
	if usr.Activated != 1 {
		//用户已经失效
		env.GenJSONResponse(context, env.LoginFailedUserInactivated, nil)
		//这里记录用户失效
		database.GetDatabase().RecordUserInActivated(key)
		//context.JSON(env.LoginFailedUserInactivated, gin.H{"status": language.GetText(env.LoginFailedUserInactivated)})
		return
	}
	loc, lok := time.LoadLocation("Local")
	regDate, pok := time.ParseInLocation(env.FullDateTimeFormat, usr.RegDate, loc)
	expDate, eok := time.ParseInLocation(env.FullDateTimeFormat, usr.ExpDate, loc)
	if lok != nil || pok != nil || eok != nil {
		//获取用户时间出错
		env.GenJSONResponse(context, env.LoginFailedGetDateFailed, nil)
		//context.JSON(env.LoginFailedGetDateFailed, gin.H{"status": language.GetText(env.LoginFailedGetDateFailed)})
		return
	}
	locExpDate := regDate.Add(time.Minute * time.Duration(usr.Duration))
	now := time.Now()
	if now.After(locExpDate) || now.After(expDate) {
		//用户已经过期
		env.GenJSONResponse(context, env.LoginFailedUserExpired, nil)
		//这里记录用户过期
		database.GetDatabase().RecordUserExpired(key)
		//context.JSON(env.LoginFailedUserExpired, gin.H{"status": language.GetText(env.LoginFailedUserExpired)})
		return
	}

	deltaDuration := locExpDate.Sub(now)
	duration := int64(deltaDuration.Minutes())

	//生成TOKEN
	token := crypto.GenToken(duration, usr.Salt)

	if strings.Compare(token, "") == 0 {
		env.GenJSONResponse(context, env.LoginFailedGenTokenFailed, nil)
		//context.JSON(env.LoginFailedGenTokenFailed, gin.H{"status": language.GetText(env.LoginFailedGenTokenFailed)})
		return
	}

	//这里记录登录成功次数
	database.GetDatabase().RecordUserLoginSuccess(key)
	//这里记录登录IP
	//database.GetDatabase().RecordUserLoginIP(key, "IP")
	//客户端需要维护token
	env.GenJSONResponseWithMsg(context, env.LoginSuccess, token)
	//context.JSON(env.LoginSuccess, gin.H{"status": token})
}

//LogoutWechat 登出控制器  暂时无用
type LogoutWeChat struct {
	Controllers
}

//GetGroup 空
func (p *LogoutWeChat) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 /usr/logout
func (p *LogoutWeChat) GetRelativePath() string {
	return "/usr/logoutwechat"
}

//GetMethod 方法 POST
func (p *LogoutWeChat) GetMethod() string {
	return "POST"
}

//GetFunc 注销方法实现
func (p *LogoutWeChat) GetFunc() func(context *gin.Context) {
	return p.logoutwechat
}

//GetFunc 注销方法实现
func (p *LogoutWeChat) logoutwechat(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : logoutwechat() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("logout %+v", context))
	// Parse JSON
	logoutRequest := models.LogoutRequest{}
	if context.Bind(&logoutRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(env.LogoutFailedParamsErrors, gin.H{"status": "登出失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s", logoutRequest.User, logoutRequest.P1)
	usr := database.GetDatabase().GetUserByName(logoutRequest.User)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.LogoutFailedUserDoNotExists, nil)
		//context.JSON(env.LogoutFailedUserDoNotExists, gin.H{"status": "登出失败，用户不存在"})
		return
	}
	//登出成功
	token := crypto.GenToken(0, usr.Salt)
	if strings.Compare(token, "") == 0 {
		env.GenJSONResponse(context, env.LogoutFailed, nil)
		//context.JSON(env.LogoutFailed, gin.H{"status": "登出失败"})
		return
	}
	//这里记录登出时间
	database.GetDatabase().RecordUserLogoutSuccess(key)
	env.GenJSONResponse(context, env.LogoutSuccess, nil)
	//context.JSON(env.LogoutSuccess, gin.H{"status": "登出成功?"})
}

//UserInfoWeChat 用户信息控制器
type UserInfoWeChat struct {
	Controllers
}

//GetGroup 空
func (p *UserInfoWeChat) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 /usr/logout
func (p *UserInfoWeChat) GetRelativePath() string {
	return "/usr/infowechat"
}

//GetMethod 方法 POST
func (p *UserInfoWeChat) GetMethod() string {
	return "POST"
}

//GetFunc 获取用户信息方法实现
func (p *UserInfoWeChat) GetFunc() func(context *gin.Context) {
	return p.getInfo
}

//GetFunc 获取用户信息方法实现
func (p *UserInfoWeChat) getInfo(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("users.go : getInfowechat() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("getInfo %+v", context))
	// Parse JSON
	userInfoRequest := models.UserInfoRequest{}
	if context.Bind(&userInfoRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		//context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，参数不正确"})
		return
	}
	key := fmt.Sprintf("%s%s", userInfoRequest.User, userInfoRequest.P1)
	log.Info(key)
	usr := database.GetDatabase().GetUserByName(userInfoRequest.User)
	//查询用户是否存在
	if usr == nil {
		//用户不存在 返回
		env.GenJSONResponse(context, env.GetUserInfoFailedUserDoNotExists, nil)
		//context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"status": "获取用户信息失败，用户不存在"})
		return
	}

	//创建回报
	rsp := models.UserInfoResponse{
		User:      usr.User,
		FullName:  usr.FullName,
		Company:   usr.Company,
		Address:   usr.Address,
		Email:     usr.Email,
		Phone:     usr.Phone,
		Space:     usr.Space,
		Tables:    usr.Tables,
		Functions: usr.Functions,
		RegDate:   usr.RegDate,
		ExpDate:   usr.ExpDate,
	}
	info, err := json.Marshal(rsp)
	if err != nil {
		info = []byte("{}")
		env.GenJSONResponse(context, env.GetUserInfoFailed, nil)
		return
	}
	env.GenJSONResponseWithMsg(context, env.GetUserInfoSuccess, string(info))
	//info := fmt.Sprintf("{\"用户名\":\"%s\",\"全名\":\"%s\",\"公司\":\"%s\",\"地址\":\"%s\",\"电话\":\"%s\",\"邮箱\":\"%s\",\"注册\":\"%s\",\"过期\":\"%s\"}", usr.User, usr.FullName, usr.Company, usr.Address, usr.Phone, usr.Email, usr.RegDate, usr.ExpDate)
	//context.JSON(http.StatusOK, gin.H{"status": "获取用户信息成功", "info": string(info)})
}
