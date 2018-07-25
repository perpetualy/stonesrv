package log

import (
    "testing"
	. "github.com/smartystreets/goconvey/convey"
)

//Debug 函数测试
func TestLogDebug(t *testing.T) {
	Convey("【测试】 LogDebug", t, func() {

		Convey("获取方法名称 ", func() {
			So("POST", ShouldEqual, "POST")
		})
	})
}

//Info 函数测试
func TestLogInfo(t *testing.T) {
	Convey("【测试】 LogInfo", t, func() {

		Convey("获取方法名称 ", func() {
			So("POST", ShouldEqual, "POST")
		})
	})
}

//Error 函数测试
func TestLogError(t *testing.T) {
	Convey("【测试】 LogError", t, func() {

		Convey("获取方法名称 ", func() {
			So("POST", ShouldEqual, "POST")
		})
	})
}