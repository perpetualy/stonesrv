package env

import (
	"fmt"
	"xmvideo/language"

	"github.com/gin-gonic/gin"
)

//GenJSONResponse 生成消息回复，可以带1个格式化参数
func GenJSONResponse(context *gin.Context, code int, param interface{}) {
	if param == nil {
		context.JSON(code, gin.H{"status": language.GetText(code)})
	} else {
		context.JSON(code, gin.H{"status": fmt.Sprintf(language.GetText(code), param)})
	}
}

//GenJSONResponseWithMsg 生成消息回复，带固定字符串
func GenJSONResponseWithMsg(context *gin.Context, code int, msg string) {
	context.JSON(code, gin.H{"status": msg})
}

//GetCodeText 获取编码的文字
func GetCodeText(code int) string {
	return language.GetText(code)
}
