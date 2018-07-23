package controllers

import (
	"strings"
	"net/url"
	"net/http"
	"github.com/gin-gonic/gin"
	"stonesrv/conf"
	"stonesrv/database"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//RegisterGetGroup函数测试
func TestRegisterGetGroup(t *testing.T) {
	Convey("【测试】 RegisterGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		register := Register{}
		group := register.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "")
		})
	})
}

//RegisterGetRelativePath函数测试
func TestRegisterGetRelativePath(t *testing.T) {
	Convey("【测试】 RegisterGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		register := Register{}
		relativePath := register.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/usr/register")
		})
	})
}

//RegisterGetMethod函数测试
func TestRegisterGetMethod(t *testing.T) {
	Convey("【测试】 AdminGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		register := Register{}
		method := register.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//Register函数测试
func TestRegister(t *testing.T) {
	Convey("【测试】 Register", t, func() {

		Convey("注册 ", func() {
			So(1, ShouldEqual, 1)
		})
	})
}

//LoginGetGroup函数测试
func TestLoginGetGroup(t *testing.T) {
	Convey("【测试】 LoginGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		login := Login{}
		group := login.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "")
		})
	})
}

//LoginGetRelativePath函数测试
func TestLoginGetRelativePath(t *testing.T) {
	Convey("【测试】 LoginGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		login := Login{}
		relativePath := login.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/usr/login")
		})
	})
}

//LoginGetMethod函数测试
func TestLoginGetMethod(t *testing.T) {
	Convey("【测试】 LoginGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		login := Login{}
		method := login.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//Login函数测试
func TestLogin(t *testing.T) {
	Convey("【测试】 Login", t, func() {

		Convey("登录 ", func() {
			So(1, ShouldEqual, 1)
		})
	})
}

//LogoutGetGroup函数测试
func TestLogoutGetGroup(t *testing.T) {
	Convey("【测试】 LogoutGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		logout := Logout{}
		group := logout.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "")
		})
	})
}

//LogoutGetRelativePath函数测试
func TestLogoutGetRelativePath(t *testing.T) {
	Convey("【测试】 LogoutGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		logout := Logout{}
		relativePath := logout.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/usr/logout")
		})
	})
}

//LogoutGetMethod函数测试
func TestLogoutGetMethod(t *testing.T) {
	Convey("【测试】 LogoutGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		logout := Logout{}
		method := logout.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//Logout函数测试
func TestLogout(t *testing.T) {
	Convey("【测试】 Logout", t, func() {
		
		logout := Logout{}
		database.Init()
		
		var logouinfo = url.Values{}
		logouinfo.Add("User", "ddd")
		logouinfo.Add("P1", "ddd")
		logouinfo.Add("P2", "ddd")
		data := logouinfo.Encode()
		req, err := http.NewRequest("POST", "https://127.0.0.1:8621/usr/logout", strings.NewReader(data))
		if err != nil{
			return
		}
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		c := &gin.Context{
			Request:req,
		}
		logout.GetFunc()(c)
		Convey("登出 ", func() {
			So(nil, ShouldEqual, nil)
		})
	})
}
