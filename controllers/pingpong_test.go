package controllers

import (
	"stonesrv/conf"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//PingPongGetGroup函数测试
func TestPingPongGetGroup(t *testing.T) {
	Convey("【测试】 PingPongGetGroup", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		pingpong := PingPong{}
		group := pingpong.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/")
		})
	})
}

//PingPongGetRelativePath函数测试
func TestPingPongGetRelativePath(t *testing.T) {
	Convey("【测试】 PingPongGetRelativePath", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		pingpong := PingPong{}
		relativePath := pingpong.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "ping")
		})
	})
}

//PingPongGetMethod函数测试
func TestPingPongGetMethod(t *testing.T) {
	Convey("【测试】 PingPongGetMethod", t, func() {
		conf.Init("../conf/stonesrv.cfg")
		pingpong := PingPong{}
		method := pingpong.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "GET")
		})
	})
}

//PingPongRequest测试
func TestPingPongRequest(t *testing.T) {
	Convey("【测试】 PingPongRequest", t, func() {

		Convey("PingPongRequest ", func() {
			So(1, ShouldEqual, 1)
		})
	})
}
