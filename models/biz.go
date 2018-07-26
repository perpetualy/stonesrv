package models

//LoginRequest 登录请求
type LoginRequest struct {
	User      	string `json:"user" binding:"required"`		//用户名
	Password  	string `json:"password" binding:"required"` 	//MD5 + SHA1 + SALT
	P1       	string `json:"p1" binding:"required"`			//MAC MD5
	P2     		string `json:"p2" binding:"required"`			//Disk0 MD5
}

//LogoutRequest 登出请求
type LogoutRequest struct {
	User      	string `json:"user" binding:"required"`		//用户名
	P1       	string `json:"p1" binding:"required"`			//MAC MD5
	P2     		string `json:"p2" binding:"required"`			//Disk0 MD5
}

//UserInfoRequest 获取用户信息请求
type UserInfoRequest struct {
	User      	string `json:"user" binding:"required"`		//用户名
	P1       	string `json:"p1" binding:"required"`			//MAC MD5
	P2     		string `json:"p2" binding:"required"`			//Disk0 MD5
}

//UserInfoResponse 获取用户信息回报
type UserInfoResponse struct {
	User      	string `json:"user" binding:"required"`		//用户名
	FullName    string `json:"fullname" binding:"required"`		//全名
	Company     string `json:"company" binding:"required"`		//公司名
	Address		string `json:"address" binding:"required"`		//地址
	Phone		string `json:"phone" binding:"required"`		//电话
	Email		string `json:"email" binding:"required"`		//邮箱
	RegDate		string `json:"regdate" binding:"required"`		//注册日期
	ExpDate		string `json:"expdate" binding:"required"`		//到期日
}

//UpdateRequest 版本更新请求
type UpdateRequest struct {
	Version 	string `json:"version" binding:"required"`		//本地版本号
	MD5	    	string `json:"md5" binding:"required"`			//本地文件MD5
}

//UpdateResponse 版本更新回报
type UpdateResponse struct{
	Version 		string `json:"version"`						//服务器版本号
	MD5				string `json:"md5"`							//服务器文件MD5
	Info    		string `json:"info"`						//更新内容
	RelDate 		string `json:"reldate"`						//更新日期
}