package controllers

import (
	"stonesrv/conf"
	"github.com/gin-gonic/gin"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

//AdminGetGroup函数测试
func TestAdminGetGroup(t *testing.T) {
	Convey("【测试】 AdminGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		admin := Admin{}
		group := admin.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/")
		})
	})
}

//AdminGetRelativePath函数测试
func TestAdminGetRelativePath(t *testing.T) {
	Convey("【测试】 AdminGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		admin := Admin{}
		relativePath := admin.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "admin")
		})
	})
}

//AdminGetMethod函数测试
func TestAdminGetMethod(t *testing.T) {
	Convey("【测试】 AdminGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		admin := Admin{}
		method := admin.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//Adminlogin函数测试
func TestAdminLogin(t *testing.T) {
	Convey("【测试】 AdminLogin", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		admin := &Admin{}
		c := new(gin.Context)
		c.Set("user", "abcd")
		c.Set("Value", "1234")
		c.Set("MyData", "456789")

		Convey("登录 ", func() {
			admin.login(c)
			//So(method, ShouldEqual, "POST")
		})
	})
}