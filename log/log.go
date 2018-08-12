package log

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

var l = new(Log)

//Log 日志结构体
type Log struct {
	fileLog       *logrus.Logger
	stdLog        *logrus.Logger
	bAdditionInfo bool
	bConsole      bool
}

//Init 初始化日志
func Init(path string, console bool, addition bool) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	l.fileLog = logrus.New()
	l.stdLog = logrus.New()
	l.bConsole = console
	l.bAdditionInfo = addition

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		l.fileLog.Out = file
	}
	l.stdLog.Out = os.Stdout
}

//Debug 调试日志
func Debug(msg string) {
	ouput := msg
	if l.bAdditionInfo {
		ouput = fmt.Sprintf("%s [%s]", msg, getAdditionInfo())
	}

	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Debug]",
		}).Debug(ouput)
	}

	if l.bConsole {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Debug]",
		}).Debug(ouput)
	}

}

//Info 一般日志
func Info(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Info]",
		}).Info(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Info]",
		}).Info(msg)
	}

}

//Warn 警告
func Warn(msg string) {
	ouput := msg
	if l.bAdditionInfo {
		ouput = fmt.Sprintf("%s [%s]", msg, getAdditionInfo())
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Warn!]",
		}).Warn(ouput)
	}
	if l.bConsole {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Warn!]",
		}).Warn(ouput)
	}
}

//Error 错误日志
func Error(msg string) {
	ouput := msg
	if l.bAdditionInfo {
		ouput = fmt.Sprintf("%s [%s]", msg, getAdditionInfo())
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Error!!!]",
		}).Error(ouput)
	}
	if l.bConsole {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Error!!!]",
		}).Error(ouput)
	}
}

//Fatal 错误后退出
func Fatal(msg string) {
	ouput := msg
	if l.bAdditionInfo {
		ouput = fmt.Sprintf("%s [%s]", msg, getAdditionInfo())
	}

	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[~~~Fatal~~~]",
		}).Fatal(ouput)
	}
	if l.bConsole {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[~~~Fatal~~~]",
		}).Fatal(ouput)
	}
}

//Panic 错误后panic
func Panic(msg string) {
	ouput := msg
	if l.bAdditionInfo {
		ouput = fmt.Sprintf("%s [%s]", msg, getAdditionInfo())
	}

	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[~~~!!!Panic!!!~~~]",
		}).Panic(ouput)
	}
	if l.bConsole {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[~~~!!!Panic!!!~~~]",
		}).Panic(ouput)
	}
}

func getAdditionInfo() string {
	funcName, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("Func Name : %s, File : %s, Line=%d", runtime.FuncForPC(funcName).Name(), file, line)
	}
	return "No more addition"
}
