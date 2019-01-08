package controllers

import (
	"fmt"
	"runtime/debug"
	"stonesrv/database"
	"stonesrv/env"
	"stonesrv/log"
	"stonesrv/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
)

//DailyInfo 每日一句用户
type DailyInfo struct {
	Controllers
}

//GetGroup 空
func (p *DailyInfo) GetGroup() string {
	return ""
}

//GetRelativePath 路径/usr/register
func (p *DailyInfo) GetRelativePath() string {
	return "/dailyinfo"
}

//GetMethod 方法 GET
func (p *DailyInfo) GetMethod() string {
	return "GET"
}

//GetFunc 注册方法实现
func (p *DailyInfo) GetFunc() func(context *gin.Context) {
	return p.dailyinfo
}

func (p *DailyInfo) dailyinfo(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("dailyinfo.go : dailyinfo() %v %+v", err, string(debug.Stack())))
		}
	}()

	req, ok := context.GetQuery("req")
	if !ok {
		env.GenJSONResponse(context, env.ParamsErrors, nil)
		return
	}
	getdailyinfoRequest := models.GetDailyInfoRequest{
		Req: req,
	}
	log.Info(fmt.Sprintf("Req:%v", getdailyinfoRequest))

	dailyInfo := database.GetDatabase().GetDailyInfo()
	if dailyInfo == nil {
		env.GenJSONResponse(context, env.GetDailyInfoFailed, nil)
		return
	}
	info, err := json.Marshal(dailyInfo)
	if err != nil {
		info = []byte("{}")
		env.GenJSONResponse(context, env.GetDailyInfoFailed, nil)
		return
	}
	env.GenJSONResponseWithMsg(context, env.GetDailyInfoSuccess, string(info))

}
