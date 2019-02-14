package controllers

import (
	"fmt"
	"runtime/debug"
	"time"
	"xmvideo/crypto"
	"xmvideo/database"
	"xmvideo/env"
	"xmvideo/log"
	"xmvideo/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
)

//InsertPack 插入套餐控制器
type InsertPack struct {
	Controllers
}

//GetGroup 空
func (p *InsertPack) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 /auth/insertpack
func (p *InsertPack) GetRelativePath() string {
	return "/insertpack"
}

//GetMethod 方法 POST
func (p *InsertPack) GetMethod() string {
	return "POST"
}

//GetFunc 插入套餐方法实现
func (p *InsertPack) GetFunc() func(context *gin.Context) {
	return p.insertPack
}

//insertPack 插入套餐方法实现
func (p *InsertPack) insertPack(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("packs.go : insertPack() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("insert pack %+v", context))
	// Parse JSON
	insertPackRequest := models.InsertPackRequest{}
	if context.Bind(&insertPackRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		return
	}

	pack := database.GetDatabase().GetPack(insertPackRequest.User)
	//查询PACK是否存在 如果存在PACK 不干啥
	if pack != nil {
		log.Info(fmt.Sprintf("last pack is %+v", pack))
	}
	user := database.GetDatabase().GetUserByName(insertPackRequest.User)
	if user == nil {
		env.GenJSONResponse(context, env.InsertPackFailedGetUser, nil)
		return
	}

	key := fmt.Sprintf("%s_%s_%s", user.Key, insertPackRequest.WeChatID, insertPackRequest.OrderID)
	if !database.GetDatabase().IsUserPaied(key) {
		env.GenJSONResponse(context, env.InsertPackFailedNoPaied, nil)
		return
	}

	now := time.Now()
	createTime := now.Format(env.FullDateTimeFormat)
	locExpTime := now.Add(time.Minute * time.Duration(insertPackRequest.Duration))
	expTime := locExpTime.Format(env.FullDateTimeFormat)
	token := crypto.GenToken(insertPackRequest.Duration, int64(crypto.PackSalt))
	userPack := models.UserPack{
		PackKey:    insertPackRequest.PackKey,
		Name:       insertPackRequest.Name,
		Desc:       insertPackRequest.Desc,
		UserKey:    user.Key,
		Space:      insertPackRequest.Space,
		Tables:     insertPackRequest.Tables,
		Functions:  insertPackRequest.Functions,
		CreateTime: createTime,
		ExpTime:    expTime,
		Token:      token,
	}
	err := database.GetDatabase().InsertPack(userPack)
	if err != nil {
		env.GenJSONResponse(context, env.InsertPackFailed, nil)
		return
	}
	env.GenJSONResponse(context, env.InsertPackSuccess, nil)
}

//GetPackToken 获取套餐TOKEN控制器
type GetPackToken struct {
	Controllers
}

//GetGroup 空
func (p *GetPackToken) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 /auth/getpacktoken
func (p *GetPackToken) GetRelativePath() string {
	return "/getpacktoken"
}

//GetMethod 方法 POST
func (p *GetPackToken) GetMethod() string {
	return "POST"
}

//GetFunc 获取套餐TOKEN方法实现
func (p *GetPackToken) GetFunc() func(context *gin.Context) {
	return p.getPackToken
}

//getPackToken 获取套餐TOKEN方法实现
func (p *GetPackToken) getPackToken(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("packs.go : GetPackToken() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("get pack token %+v", context))
	// Parse JSON
	getPackTokenRequest := models.GetPackTokenRequest{}
	if context.Bind(&getPackTokenRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		return
	}

	pack := database.GetDatabase().GetPack(getPackTokenRequest.User)
	//查询PACK是否存在 如果不存在PACK 退出
	if pack == nil {
		env.GenJSONResponse(context, env.GetPackTokenFailedNoPack, nil)
		return
	}

	//返回PACK TOKEN
	env.GenJSONResponseWithMsg(context, env.GetPackTokenSuccess, pack.Token)
}

//GetPack 获取套餐控制器
type GetPack struct {
	Controllers
}

//GetGroup 空
func (p *GetPack) GetGroup() string {
	return "/pack"
}

//GetRelativePath 路径 /pack/getpack
func (p *GetPack) GetRelativePath() string {
	return "/getpack"
}

//GetMethod 方法 POST
func (p *GetPack) GetMethod() string {
	return "POST"
}

//GetFunc 获取套餐方法实现
func (p *GetPack) GetFunc() func(context *gin.Context) {
	return p.getPack
}

//getPack 获取套餐方法实现
func (p *GetPack) getPack(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("packs.go : GetPack() %v %+v", err, string(debug.Stack())))
		}
	}()
	//log.Info(fmt.Sprintf("get pack %+v", context))
	// Parse JSON
	getPackRequest := models.GetPackRequest{}
	if context.Bind(&getPackRequest) != nil {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		return
	}

	pack := database.GetDatabase().GetPack(getPackRequest.User)
	//查询PACK是否存在 如果不存在PACK 退出
	if pack == nil {
		env.GenJSONResponse(context, env.GetPackFailedNoPack, nil)
		return
	}

	pack.Token = ""
	pack.Key = ""
	pack.UserKey = ""
	v, _ := json.Marshal(pack)
	//返回PACK
	env.GenJSONResponseWithMsg(context, env.GetPackSuccess, string(v))
}
