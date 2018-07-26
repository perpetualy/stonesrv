package models

//Config ini配置文件
type Config struct {
	DBAddress     	string `ini:"DBAddress"`
	ServerAddress 	string `ini:"ServerAddress"`
	ServerPort    	string `ini:"ServerPort"`
	SSLCrtFile    	string `ini:"SSLCrtFile"`
	SSLKeyFile    	string `ini:"SSLKeyFile"`
}

//User 用户信息
type User struct {
	Key      		string `json:"_key,omitempty"`
	User     		string `json:"user" binding:"required"`
	Password 		string `json:"password" binding:"required"`
	Email    		string `json:"email" binding:"required"`
	Address  		string `json:"address" binding:"required"`
	FullName 		string `json:"fullname" binding:"required"`
	Phone    		string `json:"phone" binding:"required"`
	Company  		string `json:"company" binding:"required"`

	Mac       		string `json:"mac" binding:"required"`
	Disk0     		string `json:"disk0" binding:"required"`
	Salt      		int64  `json:"salt" binding:"required"`		//暂时无用 强制设定为2
	Duration  		int64  `json:"duration" binding:"required"`
	RegDate   		string `json:"regdate"`
	ExpDate   		string `json:"expdate"`
	Activated 		int64  `json:"activated"`
}

//MAC 注册的MAC
type MAC struct {
	Key     		string `json:"_key,omitempty"`
	UserKey 		string `json:"userkey"`
}

//Disk0 注册的Disk0
type Disk0 struct {
	Key     		string `json:"_key,omitempty"`
	UserKey 		string `json:"userkey"`
}

//Update 版本更新
type Update struct {
	Key     		string `json:"_key,omitempty"`
	Version 		string `json:"version"`
	MD5				string `json:"md5"`
	Info    		string `json:"info"`
	Force			int64  `json:"force"`
	RelDate 		string `json:"reldate"`
}
