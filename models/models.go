//Package models models.go 数据库对象   [Author:FlynnYal CreateTime:2018-07]
package models

import (
	"encoding/xml"
)

//Config ini配置文件
type Config struct {
	DBAddress     string `ini:"DBAddress"`
	DBUser        string `ini:"DBUser"`
	DBPassword    string `ini:"DBPassword"`
	DBName        string `ini:"DBName"`
	ServerAddress string `ini:"ServerAddress"`
	ServerPort    string `ini:"ServerPort"`
	SSLCrtFile    string `ini:"SSLCrtFile"`
	SSLKeyFile    string `ini:"SSLKeyFile"`
	UpdatesDir    string `ini:"UpdatesDir"`
	UpdateFile    string `ini:"UpdateFile"`
	Language      string `ini:"Language"`
	DebugMode     string `ini:"DebugMode"`
	ToCMode       string `ini:"ToCMode"`
}

//Word 多语言基本单词
type Word struct {
	XMLName xml.Name `xml:"Word"`
	Code    int      `xml:"Code"`
	Text    string   `xml:"Text"`
}

//Words 多语言XML字符列表
type Words struct {
	XMLName xml.Name `xml:"Words"`
	Word    []Word   `xml:"Word"`
}

//User 用户信息
type User struct {
	Key      string `json:"_key,omitempty"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Address  string `json:"address" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Company  string `json:"company" binding:"required"`

	Space     int64 `json:"space" binding:"required"`     //用户空间限制
	Tables    int64 `json:"tables" binding:"required"`    //用户表限制
	Functions int64 `json:"functions" binding:"required"` //用户功能限制

	Mac       string `json:"mac" binding:"required"`
	Disk0     string `json:"disk0" binding:"required"`
	Salt      int64  `json:"salt" binding:"required"` //暂时无用 强制设定为2
	Duration  int64  `json:"duration" binding:"required"`
	RegDate   string `json:"regdate"`
	ExpDate   string `json:"expdate"`
	Activated int64  `json:"activated"`
}

//MAC 注册的MAC
type MAC struct {
	Key     string `json:"_key,omitempty"`
	UserKey string `json:"userkey"`
}

//Disk0 注册的Disk0
type Disk0 struct {
	Key     string `json:"_key,omitempty"`
	UserKey string `json:"userkey"`
}

//UserBehavior 用户行为
type UserBehavior struct {
	Key            string   `json:"_key,omitempty"`
	PasswordFailed int64    `json:"passwordfailed"`
	InActivated    int64    `json:"inactivated"`
	Expired        int64    `json:"expired"`
	LoginSuccess   int64    `json:"loginsuccess"`
	UsedSpace      int64    `json:"usedspaces"`
	UsedTable      int64    `json:"usedtables"`
	LastLogin      string   `json:"lastlogin"`
	LastLogout     string   `json:"lastlogout"`
	LoginIPs       []string `json:"loginips"`
}

//UserPack 用户套餐
type UserPack struct {
	Key        string `json:"_key,omitempty"` //自动生成
	Name       string `json:"name"`           //名称
	Desc       string `json:"desc"`           //描述
	UserKey    string `json:"userkey"`        //用户KEY
	Space      int64  `json:"space"`          //用户空间限制
	Tables     int64  `json:"tables"`         //用户表限制
	Functions  int64  `json:"functions"`      //用户功能限制
	CreateTime string `json:"createtime"`     //创建时间
	ExpTime    string `json:"exptime"`        //过期时间
	Token      string `json:"token"`          //token
}

//UserSpacePlus 用户空间叠加包
type UserSpacePlus struct {
	Key        string `json:"_key,omitempty"` //自动生成
	Name       string `json:"name"`           //名称
	Desc       string `json:"desc"`           //描述
	UserKey    string `json:"userkey"`        //用户KEY
	Space      int64  `json:"space"`          //用户空间限制
	CreateTime string `json:"createtime"`     //创建时间
	ExpTime    string `json:"exptime"`        //过期时间
	Token      string `json:"token"`          //token
}

//UserTablePlus 用户表叠加包
type UserTablePlus struct {
	Key        string `json:"_key,omitempty"` //自动生成
	Name       string `json:"name"`           //名称
	Desc       string `json:"desc"`           //描述
	UserKey    string `json:"userkey"`        //用户KEY
	Tables     int64  `json:"tables"`         //用户表限制
	CreateTime string `json:"createtime"`     //创建时间
	ExpTime    string `json:"exptime"`        //过期时间
	Token      string `json:"token"`          //token
}

//UserSpaceAndTablePlus 用户空间和表叠加包
type UserSpaceAndTablePlus struct {
	Key        string `json:"_key,omitempty"` //自动生成
	Name       string `json:"name"`           //名称
	Desc       string `json:"desc"`           //描述
	UserKey    string `json:"userkey"`        //用户KEY
	Space      int64  `json:"space"`          //用户空间限制
	Tables     int64  `json:"tables"`         //用户表限制
	CreateTime string `json:"createtime"`     //创建时间
	ExpTime    string `json:"exptime"`        //过期时间
	Token      string `json:"token"`          //token
}

//UserOrderWeChat 用户订单微信关系表
type UserOrderWeChat struct {
	Key        string `json:"_key,omitempty"`              //自动生成
	UserKey    string `json:"userkey" binding:"required"`  //用户KEY
	OrderID    string `json:"orderid" binding:"required"`  //订单号
	WeChatID   string `json:"wechatid" binding:"required"` //关联的微信号
	CreateTime string `json:"createtime"`                  //创建时间
	UpdateTime string `json:"updatetime"`                  //更新时间
	Status     string `json:"status"`                      //订单状态
}

//Updates 版本更新
type Updates struct {
	Key     string `json:"_key,omitempty"`
	Version string `json:"version"`
	MD5     string `json:"md5"`
	Info    string `json:"info"`
	Force   int64  `json:"force"`
	RelDate string `json:"reldate"`
}
