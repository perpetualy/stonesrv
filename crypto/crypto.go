package crypto

import (
	"time"
	"crypto/md5"
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var crypt = initCrypto()

//Crypto 加密算法
type Crypto struct {
	secrctSeeds	[]string
}

//Md5 MD5加密
func Md5(rawstring string) string{
	return crypt.makeMd5(rawstring)
}

//Bcrypt Bcrypt加密
func Bcrypt(rawstring string) string{
	return crypt.makeBcrypt(rawstring)
}

//SHA512 SHA512加密
func SHA512(rawstring string) string{
	return crypt.makeSHA512(rawstring)
}

//SHA256 SHA256加密
func SHA256(rawstring string) string{
	return crypt.makeSHA256(rawstring)
}

//GenToken 生成TOKEN
func GenToken(duration int64, salt int64) string{
	l := len(crypt.secrctSeeds)
	ind := int(salt) % l
	return crypt.genToken(duration, ind)
}

//GetSecrct 获取SECRCT
func GetSecrct(salt int) []byte{
	l := len(crypt.secrctSeeds)
	ind := salt % l
	return []byte(crypt.secrctSeeds[ind])
}

func initCrypto() *Crypto {
	seeds := []string{"LUU)-DR/=vox?`\\WD(HE8:$L;_lwW1;)r7D4&K/\\@VObWCEt'qQ\\57TD1JTa7o1i","Ig_c+/*ayS\\ocMP'.U>G~t\"ALV>+{JwM\"U?iWpQ8yvX:z\\7|9}lpfa.qkGjTDBvS","D0z!{\\#[d2KNtx8_[Eq'S5FsdAoivD1%*s|>ie)S^8KrIgnM\\tWlo&DaaCr{6]U'"}
	c := &Crypto{
		secrctSeeds: seeds,
	}
	return c
}

//MakeMd5 用MD5加密后的字符串	DICK0 和 MAC用MD5够了
func (p *Crypto) makeMd5(rawstring string) string {
	md := md5.New()
	return string(md.Sum([]byte(rawstring)))
}

//MakeBcrypt 用Bcrypt加密后的字符串
func (p *Crypto) makeBcrypt(rawstring string) string {
	bcryptstr, err := bcrypt.GenerateFromPassword([]byte(rawstring), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bcryptstr)
}

//MakeSHA512 用SHA512加密后的字符串
func (p *Crypto) makeSHA512(rawstring string) string {
	sh := sha512.New()
	return string(sh.Sum([]byte(rawstring)))
}

//MakeSHA256 用SHA256加密后的字符串
func (p *Crypto) makeSHA256(rawstring string) string {
	sh := sha512.New512_256()
	return string(sh.Sum([]byte(rawstring)))
}

//MakePasword 密码加密 MD5 + SHA384 + BCRYPT
//这里只是展示客户端密码如何生成
func (p *Crypto) makePasword(password string) string {
	md := md5.New()
	sh := sha512.New() //只用到512
	bcryptstr, err := bcrypt.GenerateFromPassword(sh.Sum(md.Sum([]byte(password))), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bcryptstr)
}

//GenToken 生成TOKEN
func (p *Crypto) genToken(duration int64, salt int) string {
    token := jwt.New(jwt.SigningMethodHS512)
    claims := make(jwt.MapClaims)
    claims["exp"] = time.Now().Add(time.Minute * time.Duration(duration)).Unix()
    claims["iat"] = time.Now().Unix()
    token.Claims = claims
    tokenString, err := token.SignedString([]byte(p.secrctSeeds[salt]))
    if err != nil {
        return ""
    }
    return tokenString
}
