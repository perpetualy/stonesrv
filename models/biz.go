package models

//LoginRequest 登录结构体
type LoginRequest struct {
	User      string `json:"user" binding:"required"`		//用户名
	Password  string `json:"password" binding:"required"` 	//MD5 + SHA1 + SALT
	P1       string `json:"p1" binding:"required"`			//MAC MD5
	P2     string `json:"p2" binding:"required"`			//Disk0 MD5
}

//LogoutRequest 登出结构体
type LogoutRequest struct {
	User      string `json:"user" binding:"required"`		//用户名
	P1       string `json:"p1" binding:"required"`			//MAC MD5
	P2     string `json:"p2" binding:"required"`			//Disk0 MD5
}

//UserInfoRequest 登出结构体
type UserInfoRequest struct {
	User      string `json:"user" binding:"required"`		//用户名
	P1       string `json:"p1" binding:"required"`			//MAC MD5
	P2     string `json:"p2" binding:"required"`			//Disk0 MD5
}