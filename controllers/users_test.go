package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stonesrv/conf"
	"stonesrv/database"
	"stonesrv/env"
	"stonesrv/language"
	"stonesrv/log"
	"stonesrv/models"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	configPath   = "../conf/stonesrv_test.cfg"
	languagePath = "../language"
)

func init() {
	conf.Init(configPath)
	language.Init(languagePath)
	database.Init()
}

//RegisterGetGroup函数测试
func TestRegisterGetGroup(t *testing.T) {
	Convey("【测试】 RegisterGetGroup", t, func() {
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

		register := &Register{}
		router := setupRouter(register)

		Convey("错误的方法 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/usr/register", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.PageNotFound)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldBeEmpty)
		})

		Convey("请求参数有误 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})

		user := &models.User{
			User:     "test1",
			Password: "123456",
			Email:    "abc@abc.com",
			Address:  "road1",
			FullName: "test full",
			Phone:    "22222",
			Company:  "kkkkk",
			Mac:      "33:33",
			Disk0:    "abcde",
			Salt:     2,
			Duration: 30,
		}
		Convey("新用户注册 ", func() {
			//在这之前要删除数据库中数据
			database.GetDatabase().RemoveMAC(models.MAC{Key: user.Mac})
			database.GetDatabase().RemoveDisk0(models.Disk0{Key: user.Disk0})
			database.GetDatabase().RemoveUser(models.User{Key: fmt.Sprintf("%s%s%s", user.User, user.Mac, user.Disk0)})
			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegSuccess)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))
		})
		Convey("用户已经被注册 ", func() {
			user.User = "test1"
			user.Mac = "33:33"
			user.Disk0 = "abcde"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailedUserAlreadyRegistered)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))
		})
		Convey("MAC已经被注册 ", func() {
			user.User = "test2"
			user.Mac = "33:33"
			user.Disk0 = "qwert"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailedPCAlreadyRegistered)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("Disk0已经被注册 ", func() {
			user.User = "test2"
			user.Mac = "99:99"
			user.Disk0 = "abcde"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailedPCAlreadyRegistered)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("非法的Duration ", func() {
			user.User = "test3"
			user.Mac = "12:12"
			user.Disk0 = "ghjkl"
			user.Duration = 999999

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailedInvalidDuration)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("MAC为空字符串 ", func() {
			user.User = "test2"
			user.Mac = ""
			user.Disk0 = "qwert"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("MAC为空格 ", func() {
			user.User = "test2"
			user.Mac = " "
			user.Disk0 = "qwert"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))
		})
		Convey("MAC为制表符 ", func() {
			user.User = "test2"
			user.Mac = "	"
			user.Disk0 = "qwert"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))
		})
		Convey("Disk0为空字符串 ", func() {
			user.User = "test2"
			user.Mac = "99:99"
			user.Disk0 = ""

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

		})
		Convey("Disk0为空格 ", func() {
			user.User = "test2"
			user.Mac = "99:99"
			user.Disk0 = " "

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))

		})
		Convey("Disk0为制表符 ", func() {
			user.User = "test2"
			user.Mac = "77:77"
			user.Disk0 = "	"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))

		})
		Convey("User 为空字符串 ", func() {
			user.User = ""
			user.Mac = "00:00"
			user.Disk0 = "asdfg"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

		})
		Convey("User为空格 ", func() {
			user.User = " "
			user.Mac = "01:01"
			user.Disk0 = "bnm,b"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))

		})
		Convey("User为制表符 ", func() {
			user.User = "	"
			user.Mac = "02:02"
			user.Disk0 = "345af"

			bytesData, err := json.Marshal(&user)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/register", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.RegFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code), user.User))

		})
	})
}

//LoginGetGroup函数测试
func TestLoginGetGroup(t *testing.T) {
	Convey("【测试】 LoginGetGroup", t, func() {
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
		login := &Login{}
		router := setupRouter(login)

		Convey("错误的方法 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/usr/login", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.PageNotFound)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldBeEmpty)
		})

		Convey("请求参数有误 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})

		loginReq := models.LoginRequest{
			User:     "test1",
			Password: "123456",
			P1:       "33:33",
			P2:       "abcde",
		}
		Convey("用户登录 ", func() {
			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginSuccess)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldNotBeEmpty)
			Convey(status.Status, nil)
		})
		Convey("空用户名 ", func() {
			loginReq.User = ""

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.ParamsErrors)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("空格用户名 ", func() {
			loginReq.User = " "

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserDoNotExists)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("制表符用户名 ", func() {
			loginReq.User = "	"

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserDoNotExists)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("不存在的用户名 ", func() {
			loginReq.User = "test999"

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserDoNotExists)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		Convey("密码错误 ", func() {
			loginReq.Password = "189479"

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedPasswordIncorrect)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})
		user := models.User{
			Key:       "test133:33abcde",
			User:      "test1",
			Password:  "123456",
			Email:     "abc@abc.com",
			Address:   "road1",
			FullName:  "test full",
			Phone:     "22222",
			Company:   "kkkkk",
			Mac:       "33:33",
			Disk0:     "abcde",
			Salt:      2,
			Duration:  30,
			RegDate:   time.Now().Format(env.FullDateTimeFormat),
			ExpDate:   time.Now().Add(30 * time.Minute).Format(env.FullDateTimeFormat),
			Activated: 1,
		}
		Convey("用户失效 ", func() {

			database.GetDatabase().DeactiveUser(user)
			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserInactivated)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			database.GetDatabase().ActiveUser(user)
		})
		Convey("注册时间出错 ", func() {
			tempDate := user.RegDate
			user.RegDate = "111"
			database.GetDatabase().UpdateUserInfo(user)

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedGetDateFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			user.RegDate = tempDate
			database.GetDatabase().UpdateUserInfo(user)
		})
		Convey("过期时间出错 ", func() {
			tempDate := user.ExpDate
			user.ExpDate = "111"
			database.GetDatabase().UpdateUserInfo(user)

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedGetDateFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			user.ExpDate = tempDate
			database.GetDatabase().UpdateUserInfo(user)
		})
		Convey("注册时间久远，已经过期", func() {
			tempDate := user.RegDate
			user.RegDate = "2018-01-01 00:00:00"
			database.GetDatabase().UpdateUserInfo(user)

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserExpired)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			user.RegDate = tempDate
			database.GetDatabase().UpdateUserInfo(user)
		})
		Convey("过期时间很近，已经过期", func() {
			tempDate := user.ExpDate
			user.ExpDate = time.Now().Format(env.FullDateTimeFormat)
			database.GetDatabase().UpdateUserInfo(user)

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedUserExpired)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			user.ExpDate = tempDate
			database.GetDatabase().UpdateUserInfo(user)
		})
		/* 		Convey("获取Token出错", func() {
			tempDuration := user.Duration
			user.Duration = -30
			database.GetDatabase().UpdateUserInfo(user)

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginFailedGenTokenFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			user.Duration = tempDuration
			database.GetDatabase().UpdateUserInfo(user)
		}) */
	})
}

//LogoutGetGroup函数测试
func TestLogoutGetGroup(t *testing.T) {
	Convey("【测试】 LogoutGetGroup", t, func() {
		logout := Logout{}
		group := logout.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/auth")
		})
	})
}

//LogoutGetRelativePath函数测试
func TestLogoutGetRelativePath(t *testing.T) {
	Convey("【测试】 LogoutGetRelativePath", t, func() {
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
		logout := &Logout{}
		router := setupRouter(logout)

		Convey("该用户未授权 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/usr/logout", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.AuthFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
		})

		Convey("错误的方法 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/usr/logout", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.PageNotFound)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldBeEmpty)
		})
		Convey("用户登录 ", func() {
			login := &Login{}
			routerLogin := setupRouter(login)

			loginReq := models.LoginRequest{
				User:     "test1",
				Password: "123456",
				P1:       "33:33",
				P2:       "abcde",
			}

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			routerLogin.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginSuccess)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldNotBeEmpty)
			currentToken := status.Status

			Convey("请求参数有误 ", func() {

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/logout", nil)
				req.Header.Set("Authorization", currentToken)
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.ParamsErrors)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			})

			logoutRequest := models.LogoutRequest{
				User: "test1",
				P1:   "33:33",
				P2:   "abcde",
			}
			Convey("用户登出", func() {
				bytesData, err := json.Marshal(&logoutRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/logout", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.LogoutSuccess)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			})

			Convey("用户不存在", func() {
				logoutRequest.User = "testlogout"
				bytesData, err := json.Marshal(&logoutRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/logout", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.LogoutFailedUserDoNotExists)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			})
		})

	})
}

func TestUserInfoGetGroup(t *testing.T) {
	Convey("【测试】 UserInfoGetGroup", t, func() {
		userinfo := UserInfo{}
		group := userinfo.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/auth")
		})
	})
}

//LogoutGetRelativePath函数测试
func TestUserInfoGetRelativePath(t *testing.T) {
	Convey("【测试】 UserInfoGetRelativePath", t, func() {
		userinfo := UserInfo{}
		relativePath := userinfo.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/usr/info")
		})
	})
}

//LogoutGetMethod函数测试
func TestUserInfoGetMethod(t *testing.T) {
	Convey("【测试】 UserInfoGetMethod", t, func() {
		userinfo := &UserInfo{}
		method := userinfo.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//TestUserInfo函数测试
func TestUserInfo(t *testing.T) {
	Convey("【测试】 UserInfo", t, func() {
		userinfo := &UserInfo{}
		router := setupRouter(userinfo)

		Convey("该用户未授权 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/usr/info", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.AuthFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, "该用户未授权")
		})

		Convey("错误的方法 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/usr/info", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.PageNotFound)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldBeEmpty)
		})

		Convey("用户登录 ", func() {
			login := &Login{}
			routerLogin := setupRouter(login)

			loginReq := models.LoginRequest{
				User:     "test1",
				Password: "123456",
				P1:       "33:33",
				P2:       "abcde",
			}

			bytesData, err := json.Marshal(&loginReq)
			if err != nil {
				log.Error(fmt.Sprintf("%v", err.Error()))
				return
			}
			reader := bytes.NewReader(bytesData)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/usr/login", reader)
			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			routerLogin.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.LoginSuccess)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldNotBeEmpty)
			currentToken := status.Status

			Convey("请求参数有误 ", func() {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/info", nil)
				req.Header.Set("Authorization", currentToken)
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.ParamsErrors)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			})

			userInfoRequest := models.UserInfoRequest{
				User: "test1",
				P1:   "33:33",
				P2:   "abcde",
			}

			Convey("获取用户信息", func() {
				bytesData, err := json.Marshal(&userInfoRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/info", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUserInfoSuccess)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				usr := models.UserInfoResponse{}
				json.Unmarshal([]byte(status.Status), &usr)
				So(usr.User, ShouldEqual, "test1")
			})

			Convey("用户不存在", func() {
				userInfoRequest.User = "testGetInfo"
				bytesData, err := json.Marshal(&userInfoRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/usr/info", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUserInfoFailedUserDoNotExists)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))
			})

		})

	})
}
