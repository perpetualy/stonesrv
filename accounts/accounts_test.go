package accounts

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//AddAccount函数测试
func TestAccountsAddAccount(t *testing.T) {
	Convey("【测试】 AddAccount", t, func() {
		user := AddAccount("ABC", "123")
		Convey("添加用户 ABC", func() {
			So(user, ShouldEqual, "ABC")
		})
	})
}

//GetAccounts函数测试
func TestAccountsGetAccounts(t *testing.T) {
	Convey("【测试】 GetAccounts", t, func() {
		accs := GetAccounts()
		pwd1 := accs["stone"]
		Convey("认证stone", func() {
			So(pwd1, ShouldEqual, "456789")
		})
		pwd2 := accs["STone"]
		Convey("认证STone", func() {
			So(pwd2, ShouldEqual, "")
		})
	})
}
