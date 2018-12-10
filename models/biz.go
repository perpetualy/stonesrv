//Package models biz.go 业务请求对象   [Author:FlynnYal CreateTime:2018-08]
package models

//StatusMsg 消息的内容
type StatusMsg struct {
	Status string `json:"status"` //消息具体内容
}

//Response 所有消息的返回状态
type Response struct {
	Code int       `json:"code"` //现在改为同SRV的RESPONSE CODE保持一致
	Msg  StatusMsg `json:"msg"`  //消息内容JSON
}

//LoginRequest 登录请求
type LoginRequest struct {
	User     string `json:"user" binding:"required"`     //用户名
	Password string `json:"password" binding:"required"` //MD5 + SHA1 + SALT
	P1       string `json:"p1" binding:"required"`       //MAC MD5
	P2       string `json:"p2" binding:"required"`       //Disk0 MD5
}

//LogoutRequest 登出请求
type LogoutRequest struct {
	User string `json:"user" binding:"required"` //用户名
	P1   string `json:"p1" binding:"required"`   //MAC MD5
	P2   string `json:"p2" binding:"required"`   //Disk0 MD5
}

//RegisterRequest 注册请求
type RegisterRequest struct {
	User      string `json:"user" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Address   string `json:"address" binding:"required"`
	FullName  string `json:"fullname" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Company   string `json:"company" binding:"required"`
	Space     int64  `json:"space" binding:"required"`
	Tables    int64  `json:"tables"`
	Functions int64  `json:"functions"`

	Mac       string `json:"mac" binding:"required"`
	Disk0     string `json:"disk0" binding:"required"`
	Salt      int64  `json:"salt" binding:"required"` //暂时无用 强制设定为2
	Duration  int64  `json:"duration" binding:"required"`
	RegDate   string `json:"regdate"`
	ExpDate   string `json:"expdate"`
	Activated int64  `json:"activated"`
}

//UserInfoRequest 获取用户信息请求
type UserInfoRequest struct {
	User string `json:"user" binding:"required"` //用户名
	P1   string `json:"p1" binding:"required"`   //MAC MD5
	P2   string `json:"p2" binding:"required"`   //Disk0 MD5
}

//UserInfoResponse 获取用户信息回报
type UserInfoResponse struct {
	User      string `json:"user" binding:"required"`      //用户名
	FullName  string `json:"fullname" binding:"required"`  //全名
	Company   string `json:"company" binding:"required"`   //公司名
	Address   string `json:"address" binding:"required"`   //地址
	Phone     string `json:"phone" binding:"required"`     //电话
	Email     string `json:"email" binding:"required"`     //邮箱
	Space     int64  `json:"space" binding:"required"`     //空间限制
	Tables    int64  `json:"tables" binding:"required"`    //表限制
	Functions int64  `json:"functions" binding:"required"` //功能限制
	RegDate   string `json:"regdate" binding:"required"`   //注册日期
	ExpDate   string `json:"expdate" binding:"required"`   //到期日
}

//InsertPackRequest 插入套餐请求
type InsertPackRequest struct {
	User      string `json:"user" binding:"required"`      //用户名
	Name      string `json:"name" binding:"required"`      //名称
	Desc      string `json:"desc"`                         //描述
	OrderID   string `json:"orderid" binding:"required"`   //订单号
	WeChatID  string `json:"wechatid" binding:"required"`  //关联的微信号
	Space     int64  `json:"space" binding:"required"`     //空间限制
	Tables    int64  `json:"tables" binding:"required"`    //表限制
	Functions int64  `json:"functions" binding:"required"` //功能限制
	Duration  int64  `json:"duration" binding:"required"`  //时长
}

//GetPackTokenRequest 获取套餐TOKEN请求
type GetPackTokenRequest struct {
	User string `json:"user" binding:"required"` //用户名
}

//GetPackRequest 获取套餐请求
type GetPackRequest struct {
	User string `json:"user" binding:"required"` //用户名
}

//UserBehaviorRequest 获取用户行为
type UserBehaviorRequest struct {
	User string `json:"user" binding:"required"` //用户名
}

//UpdatesRequest 版本更新请求
type UpdatesRequest struct {
	Version string `json:"version" binding:"required"` //本地版本号
	MD5     string `json:"md5" binding:"required"`     //本地文件MD5
}

//UpdatesResponse 版本更新回报
type UpdatesResponse struct {
	Version string `json:"version"` //服务器版本号
	MD5     string `json:"md5"`     //服务器文件MD5
	Info    string `json:"info"`    //更新内容
	RelDate string `json:"reldate"` //更新日期
	Path    string `json:"path"`    //下载路径
}
