package controllers

import (
	"github.com/gin-gonic/gin/json"
	"fmt"
	"strings"
	"net/http"
	"stonesrv/log"
	"stonesrv/models"
	"stonesrv/database"
	"github.com/gin-gonic/gin"
)

//Update 应用版本更新
type Update struct {
	Controllers
}

//GetGroup 空
func (p *Update) GetGroup() string {
	return "/auth"
}

//GetRelativePath 路径 update
func (p *Update) GetRelativePath() string {
	return "/update"
}

//GetMethod GET
func (p *Update) GetMethod() string {
	return "GET"
}

//GetFunc 更新方法实现
func (p *Update) GetFunc() func(context *gin.Context) {
	return func(context *gin.Context) {
		log.Info(fmt.Sprintf("login %+v", context))
		// Parse JSON
		updRequest := models.UpdateRequest{}
		if context.ShouldBind(&updRequest) != nil {
			context.JSON(http.StatusBadRequest, gin.H{"status": "获取更新失败，参数不正确"})
			return
		}
		clientVers := strings.Split(updRequest.Version, ".")
		if len(clientVers) < 3{
			context.JSON(http.StatusBadRequest, gin.H{"status": "获取更新失败，版本校验出错"})
			return
		}
		upd := database.GetDatabase().GetUpdate()
		serverVers := strings.Split(upd.Version, ".")
		if len(serverVers) < 3{
			context.JSON(http.StatusBadRequest, gin.H{"status": "没有最新版本"})
			return
		}
		if strings.Compare(upd.Version, updRequest.Version) == 0{
			context.JSON(http.StatusOK, gin.H{"status": "已经是最新版本"})
			return
		}
		if strings.Compare(upd.MD5, updRequest.MD5) == 0{
			context.JSON(http.StatusOK, gin.H{"status": "本地更新已经存在"})
			return
		}
		if strings.Compare(serverVers[0], clientVers[0]) > 0 || strings.Compare(serverVers[1], clientVers[1]) > 0 || strings.Compare(serverVers[2], clientVers[2]) > 0{
			//这里更新
			rsp := models.UpdateResponse{
				Version:upd.Version,
				MD5:upd.MD5,
				Info:upd.Info,
				RelDate:upd.RelDate,
			}
			info,err := json.Marshal(rsp)
			if err != nil{
				info = []byte("{}")
			}
			context.JSON(http.StatusOK, gin.H{"status": "发现新版本", "info": string(info)})
			return
		}
		context.JSON(http.StatusOK, gin.H{"status": "不需要更新"})
	}
}
