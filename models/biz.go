package models

//LoginRequest 登录结构体
type LoginRequest struct {
	User      string `json:"user" binding:"required"`		//用户名
	Password  string `json:"password" binding:"required"` 	//MD5 + SHA1 + SALT
	MAC       string `json:"p1" binding:"required"`			//MD5
	Disk0     string `json:"p2" binding:"required"`			//MD5
}

//LogoutRequest 登出结构体
type LogoutRequest struct {
	User      string `json:"user" binding:"required"`		//用户名
	MAC       string `json:"p1" binding:"required"`			//MD5
	Disk0     string `json:"p2" binding:"required"`			//MD5
}