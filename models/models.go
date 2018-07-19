package models

type Config struct {   //ini 配置文件
	DBAddress 		string  `ini:"DBAddress"`
}

//用户信息
type User struct {
	Key           string `json:"_key,omitempty"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Email		  string `json:"email"`
	Address 	  string `json:"address"`
	FullName	  string `json:"fullname"`

	Mac			  string `json:"mac"`
	Disk0		  string `json:"disk0"`
	Salt          int64  `json:"salt"`
	RegDate       string `json:"regdate"`
	ExpDate       string `json:"expdate"`
	Activated     int64  `json:"activated"`
}