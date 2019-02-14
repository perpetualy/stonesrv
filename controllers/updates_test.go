package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"xmvideo/conf"
	"xmvideo/database"
	"xmvideo/env"
	"xmvideo/language"
	"xmvideo/log"
	"xmvideo/models"

	. "github.com/smartystreets/goconvey/convey"
)

//UpdatesGetGroup函数测试
func TestUpdatesGetGroup(t *testing.T) {
	Convey("【测试】 UpdateGetGroup", t, func() {
		conf.Init("../conf/xmvideosrv.cfg")
		update := Updates{}
		group := update.GetGroup()
		Convey("获取分组 ", func() {
			So(group, ShouldEqual, "/auth")
		})
	})
}

//UpdatesGetRelativePath函数测试
func TestUpdatesGetRelativePath(t *testing.T) {
	Convey("【测试】 UpdateGetRelativePath", t, func() {
		conf.Init("../conf/xmvideosrv.cfg")
		update := Updates{}
		relativePath := update.GetRelativePath()
		Convey("获取相对路径 ", func() {
			So(relativePath, ShouldEqual, "/update")
		})
	})
}

//UpdatesGetMethod函数测试
func TestUpdatesGetMethod(t *testing.T) {
	Convey("【测试】 UpdateGetMethod", t, func() {
		conf.Init("../conf/xmvideosrv.cfg")
		update := Updates{}
		method := update.GetMethod()
		Convey("获取方法名称 ", func() {
			So(method, ShouldEqual, "POST")
		})
	})
}

//UpdatesRequest测试
func TestUpdatesRequest(t *testing.T) {
	Convey("【测试】 UpdateRequest", t, func() {
		updates := &Updates{}
		router := setupRouter(updates)

		Convey("该用户未授权 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/update", nil)
			router.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, env.AuthFailed)
			status := models.StatusMsg{}
			json.Unmarshal([]byte(w.Body.String()), &status)
			So(status.Status, ShouldEqual, "该用户未授权")
		})

		Convey("错误的方法 ", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/update", nil)
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
				req, _ := http.NewRequest("POST", "/auth/update", nil)
				req.Header.Set("Authorization", currentToken)
				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.ParamsErrors)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			})

			updatesRequest := models.UpdatesRequest{
				Version: "1.0.0",
				MD5:     "abcdefg",
			}

			Convey("版本校验出错 ", func() {
				updatesRequest.Version = "1.0"
				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesFailedCheckingFailed)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			})

			upd := models.Updates{
				Key:     "md5",
				Version: ".",
				MD5:     "cdefg",
				Info:    "测试版本",
				Force:   0,
				RelDate: "2018-08-08",
			}

			Convey("获取更新失败，远程没有更新 ", func() {
				database.GetDatabase().RemoveUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesFailedRemoteFailed)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

			})

			Convey("获取更新失败，远程更新版本有问题 ", func() {
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesFailedRemoteFailed)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("强制更新", func() {
				upd.Version = "0.0.1"
				upd.Force = 1
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesEmergent)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				rsp := models.UpdatesResponse{}
				json.Unmarshal([]byte(status.Status), &rsp)
				So(rsp.Version, ShouldEqual, upd.Version)
				So(rsp.Path, ShouldContainSubstring, ".zip")
				So(rsp.MD5, ShouldNotEqual, updatesRequest.MD5)
				So(rsp.Version, ShouldNotEqual, updatesRequest.Version)

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("不需要更新", func() {
				upd.Version = "1.0.0"
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesNoNeed)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("本地更新已经存在", func() {
				upd.Version = "1.0.1"
				updatesRequest.MD5 = upd.MD5
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesLocalUpdateAlready)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				So(status.Status, ShouldEqual, fmt.Sprintf(language.GetText(w.Code)))

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("发现新版本1", func() {
				upd.Version = "1.1.0"
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesUpdateFound)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				rsp := models.UpdatesResponse{}
				json.Unmarshal([]byte(status.Status), &rsp)
				So(rsp.Path, ShouldContainSubstring, ".zip")
				So(rsp.MD5, ShouldNotEqual, updatesRequest.MD5)
				So(rsp.Version, ShouldNotEqual, updatesRequest.Version)

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("发现新版本2", func() {
				upd.Version = "1.0.1"
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesUpdateFound)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				rsp := models.UpdatesResponse{}
				json.Unmarshal([]byte(status.Status), &rsp)
				So(rsp.Path, ShouldContainSubstring, ".zip")
				So(rsp.MD5, ShouldNotEqual, updatesRequest.MD5)
				So(rsp.Version, ShouldNotEqual, updatesRequest.Version)

				database.GetDatabase().RemoveUpdate(upd)
			})

			Convey("发现新版本3", func() {
				upd.Version = "1.0.0b"
				database.GetDatabase().SetUpdate(upd)

				bytesData, err := json.Marshal(&updatesRequest)
				if err != nil {
					log.Error(fmt.Sprintf("%v", err.Error()))
					return
				}
				reader := bytes.NewReader(bytesData)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/auth/update", reader)
				req.Header.Set("Authorization", currentToken)
				req.Header.Set("Content-Type", "application/json;charset=UTF-8")

				router.ServeHTTP(w, req)
				So(w.Code, ShouldEqual, env.GetUpdatesUpdateFound)
				status := models.StatusMsg{}
				json.Unmarshal([]byte(w.Body.String()), &status)
				rsp := models.UpdatesResponse{}
				json.Unmarshal([]byte(status.Status), &rsp)
				So(rsp.Path, ShouldContainSubstring, ".zip")
				So(rsp.MD5, ShouldNotEqual, updatesRequest.MD5)
				So(rsp.Version, ShouldNotEqual, updatesRequest.Version)

				database.GetDatabase().RemoveUpdate(upd)
			})
		})

	})
}
