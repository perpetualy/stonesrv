package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type Log struct {
	path string
}

//初始化日志
func Init(path string, console bool){
	f, _ := os.Create(path)
	if console{
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}else{
		gin.DefaultWriter = io.MultiWriter(f)
	}
}

//调试日志
func Debug(msg string){
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("Stone - [Debug]: %s",msg))
}

//一般日志
func Info(msg string){
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("Stone - [Info]: %s",msg))
}

//错误日志
func Error(msg string){
	fmt.Fprintln(gin.DefaultWriter, fmt.Sprintf("Stone - [Error!!!]: %s",msg))
}