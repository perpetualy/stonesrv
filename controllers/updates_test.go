package controllers

import (
	"stonesrv/conf"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//UpdatesGetGroup函数测试
func TestUpdatesGetGroup(t *testing.T) {
	Convey("【测试】 UpdateGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		group := update.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/auth")
		})
	})
}

//UpdatesGetRelativePath函数测试
func TestUpdatesGetRelativePath(t *testing.T) {
	Convey("【测试】 UpdateGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		relativePath := update.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/update")
		})
	})
}

//UpdatesGetMethod函数测试
func TestUpdatesGetMethod(t *testing.T) {
	Convey("【测试】 UpdateGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		method := update.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "GET")
		})
	})
}

//UpdatesRequest测试
func TestUpdatesRequest(t *testing.T) {
	Convey("【测试】 UpdateRequest", t, func() {

		Convey("UpdateRequest ", func() {
			So(1, ShouldEqual, 1)
		})
	})
}
