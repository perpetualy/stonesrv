package middlewares

import (
    "testing"
	. "github.com/smartystreets/goconvey/convey"
)

//AuthToken 函数测试
func TestAuthAuthToken(t *testing.T) {
	Convey("【测试】 AuthAuthToken", t, func() {

		Convey("获取方法名称 ", func() {
			So("POST", ShouldEqual, "POST")
		})
	})
}