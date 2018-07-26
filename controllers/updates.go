package controllers

import (
	"fmt"
	"net/http"
	"stonesrv/conf"
	"stonesrv/database"
	"stonesrv/log"
	"stonesrv/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
)

//Updates 应用版本更新
type Updates struct {
	Controllers
}

//GetGroup 空
func (p *Updates) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 update
func (p *Updates) GetRelativePath() string {
	return "/update"
}

//GetMethod GET
func (p *Updates) GetMethod() string {
	return "GET"
}

//GetFunc 更新方法实现
func (p *Updates) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		log.Info(fmt.Sprintf("login %+v", context))
		// Parse JSON
		updRequest := models.UpdatesRequest{}
		if context.ShouldBind(&updRequest) != nil {
			context.JSON(http.StatusBadRequest, gin.H{"status": "获取更新失败，参数不正确"})
			return
		}
		clientVers := strings.Split(updRequest.Version, ".")
		if len(clientVers) < 3 {
			context.JSON(http.StatusBadRequest, gin.H{"status": "获取更新失败，版本校验出错"})
			return
		}
		upd := database.GetDatabase().GetUpdate()
		if upd == nil {
			context.JSON(http.StatusBadRequest, gin.H{"status": "没有可用更新"})
			return
		}
		serverVers := strings.Split(upd.Version, ".")
		isForce := upd.Force == 1
		if len(serverVers) < 3 {
			context.JSON(http.StatusBadRequest, gin.H{"status": "没有最新版本"})
			return
		}
		if isForce {
			//这里强制更新
			info := p.makeResponse(upd)
			context.JSON(http.StatusOK, gin.H{"status": "紧急更新", "info": info})
			return
		}
		if strings.Compare(upd.Version, updRequest.Version) == 0 {
			context.JSON(http.StatusOK, gin.H{"status": "已经是最新版本"})
			return
		}
		if strings.Compare(upd.MD5, updRequest.MD5) == 0 {
			context.JSON(http.StatusOK, gin.H{"status": "本地更新已经存在"})
			return
		}
		if strings.Compare(serverVers[0], clientVers[0]) > 0 || strings.Compare(serverVers[1], clientVers[1]) > 0 || strings.Compare(serverVers[2], clientVers[2]) > 0 {
			//发现新版本更新
			info := p.makeResponse(upd)
			context.JSON(http.StatusOK, gin.H{"status": "发现新版本", "info": info})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": "不需要更新"})
	}
}

func (p *Updates) makeResponse(upd *models.Updates) string {
	rsp := models.UpdatesResponse{
		Version: upd.Version,
		MD5:     upd.MD5,
		Info:    upd.Info,
		RelDate: upd.RelDate,
		Path:    fmt.Sprintf("https://127.0.0.1:8621%s/%s/%s", conf.GetUpdatesDir(), upd.Version, conf.GetUpdateFile()),
	}
	info, err := json.Marshal(rsp)
	if err != nil {
		info = []byte("{}")
	}
	return string(info)
}
