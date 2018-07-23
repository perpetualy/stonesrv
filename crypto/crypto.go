package crypto

import (
	"crypto/md5"
	"golang.org/x/crypto/bcrypt"
	"crypto/sha512"
)

//MakeMd5 用MD5加密后的字符串	DICK0 和 MAC用MD5够了
func MakeMd5(rawstring string) string{
	md := md5.New()
	return string(md.Sum([]byte(rawstring)))
}

//MakeBcrypt 用Bcrypt加密后的字符串
func MakeBcrypt(rawstring string) string{
	bcryptstr, err := bcrypt.GenerateFromPassword([]byte(rawstring), bcrypt.DefaultCost)
	if err != nil{
		return ""
	}
	return string(bcryptstr)
}

//MakeSHA512 用SHA512加密后的字符串
func MakeSHA512(rawstring string) string{
	sha := sha512.New()
	return string(sha.Sum([]byte(rawstring)))
}

//MakePasword 密码加密
func MakePasword(password string) string{
	md := md5.New()
	sha := sha512.New384()	//只用到384
	bcryptstr, err := bcrypt.GenerateFromPassword(sha.Sum(md.Sum([]byte(password))), bcrypt.DefaultCost)
	if err != nil{
		return ""
	}
	return string(bcryptstr) 
}