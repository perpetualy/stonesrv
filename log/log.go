package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var l = new(Log)

//Log 日志结构体
type Log struct {
	fileLog *logrus.Logger
	stdLog  *logrus.Logger
}

//Init 初始化日志
func Init(path string, console bool) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	l.fileLog = logrus.New()
	l.stdLog = logrus.New()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		l.fileLog.Out = file
	} else {
		l.fileLog.Out = os.Stdout
	}

	if console {
		l.stdLog.Out = os.Stdout
	} else {
		l.stdLog = nil
	}
}

//Debug 调试日志
func Debug(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Debug]",
		}).Debug(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Debug]",
		}).Debug(msg)
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

//Warning 警告
func Warn(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Warn!]",
		}).Warn(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Warn!]",
		}).Warn(msg)
	}
}

//Error 错误日志
func Error(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[Error!!!]",
		}).Error(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[Error!!!]",
		}).Error(msg)
	}
}

//Fatal 错误后退出
func Fatal(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[~~~Fatal~~~]",
		}).Fatal(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[~~~Fatal~~~]",
		}).Fatal(msg)
	}
}

//Panic 错误后panic
func Panic(msg string) {
	if l.stdLog != nil {
		l.stdLog.WithFields(logrus.Fields{
			"Stone": "[~~~!!!Panic!!!~~~]",
		}).Panic(msg)
	}
	if l.fileLog != nil {
		l.fileLog.WithFields(logrus.Fields{
			"Stone": "[~~~!!!Panic!!!~~~]",
		}).Panic(msg)
	}
}
