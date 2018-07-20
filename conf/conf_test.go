package conf

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

//GetDBAddress函数测试
func TestConfGetDBAddress(t *testing.T) {
	Convey("【测试】 GetDBAddress", t, func() {
		Init("stonesrv.cfg")
		address := GetDBAddress()
		Convey("获取数据库地址 ", func() {
			So(address, ShouldEqual, "127.0.0.1")
		})
	})
}

//GetServerAddress函数测试
func TestConfGetServerAddress(t *testing.T) {
	Convey("【测试】 GetServerAddress", t, func() {
		Init("stonesrv.cfg")
		address := GetServerAddress()
		Convey("获取服务器地址 ", func() {
			So(address, ShouldEqual, "127.0.0.1")
		})

	})
}

//GetServerPort函数测试
func TestConfGetServerPort(t *testing.T) {
	Convey("【测试】 GetServerPort", t, func() {
		Init("stonesrv.cfg")
		port := GetServerPort()
		Convey("获取服务器端口 ", func() {
			So(port, ShouldEqual, "8621")
		})

	})
}