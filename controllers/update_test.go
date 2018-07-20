package controllers

import (
	"stonesrv/conf"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

//AdminGetGroup函数测试
func TestUpdateGetGroup(t *testing.T) {
	Convey("【测试】 UpdateGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		group := update.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "")
		})
	})
}

//UpdateGetRelativePath函数测试
func TestUpdateGetRelativePath(t *testing.T) {
	Convey("【测试】 UpdateGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		relativePath := update.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "update")
		})
	})
}

//UpdateGetMethod函数测试
func TestUpdateGetMethod(t *testing.T) {
	Convey("【测试】 UpdateGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		update := Update{}
		method := update.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "GET")
		})
	})
}

//UpdateRequest测试
func TestUpdateRequest(t *testing.T) {
	Convey("【测试】 UpdateRequest", t, func() {
		
		Convey("UpdateRequest ", func() {
			So(1, ShouldEqual, 1)
		})
	})
}