package crypto

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	//LoginSalt LOGIN
	LoginSalt = 2
	//PackSalt PACK
	PackSalt = 3
	//SpaceSalt SPACE PLUS
	SpaceSalt = 4
	//TableSalt TABLE PLUS
	TableSalt = 5
	//SpaceAndTableSalt SPACE AND TABLE PLUS
	SpaceAndTableSalt = 6
	//crypt 初始化加密库
	crypt = initCrypto()
)

//Crypto 加密算法
type Crypto struct {
	secrctSeeds []string
}

//Md5 MD5加密
func Md5(rawstring string) string {
	return crypt.makeMd5(rawstring)
}

//Bcrypt Bcrypt加密
func Bcrypt(rawstring string) string {
	return crypt.makeBcrypt(rawstring)
}

//SHA512 SHA512加密
func SHA512(rawstring string) string {
	return crypt.makeSHA512(rawstring)
}

//SHA256 SHA256加密
func SHA256(rawstring string) string {
	return crypt.makeSHA256(rawstring)
}

//GenToken 生成TOKEN
func GenToken(duration int64, salt int64) string {
	l := len(crypt.secrctSeeds)
	ind := int(salt) % l
	return crypt.genToken(duration, ind)
}

//GetSecrct 获取SECRCT
func GetSecrct(salt int) []byte {
	l := len(crypt.secrctSeeds)
	ind := salt % l
	return []byte(crypt.secrctSeeds[ind])
}

func initCrypto() *Crypto {
	seeds := []string{
		"LUU)-DR/=vox?`\\WD(HE8:$L;_lwW1;)r7D4&K/\\@VObWCEt'qQ\\57TD1JTa7o1i",
		"Ig_c+/*ayS\\ocMP'.U>G~t\"ALV>+{JwM\"U?iWpQ8yvX:z\\7|9}lpfa.qkGjTDBvS",
		"D0z!{\\#[d2KNtx8_[Eq'S5FsdAoivD1%*s|>ie)S^8KrIgnM\\tWlo&DaaCr{6]U'", //LOGIN
		"L86FlZW5KTEhQLM@D$l6NfD#Nwqv!449lmY4b3YS^Lyx^qeW!Xx0$SNQP9v4YlLu",   //PACK
		"T#UvFE9hYrNF1G5@ffIq5Jn6@K8%IuS$2i0KCuvuVZS@12W7Bh52X@LkowgfrR%W",   //SPACE PLUS
		"ZcPr0%N%RiFpWg*guF5^8hBx8#gR9fmf6$q*2lbJ0I%!%kHclDQwr8$ZX1syYP^W",   //TABLE PLUS
		"9yJr7S*e4Z&Ca$^*aQLLzmnlvfWoFpoZP3wkGxR@pw2vHc#cncHeKAa65Sc$SZzX",
		"xj@RhA6pnGo%kOtXx#GMMIVkoBdpKl%yFKd$9WvIlhuUCpR2B0U@c^cYa^5a^K6$",
		"maaM6P&t4%#PgjeqvWeEnY6!XB2%Y7S3XcqQuXIxfBK*a2I^2oRGzZHDgzr75Ymm",
		"#HrSYAz!9Z0$7*eo6dfd$j&75&3nvK&n7Cr2rTNJ4aoqe!BjlLRxenAC*JhlyhtO",
		"UZwTtdPCTjPipfObi8HuUb6U*Rao42t$D#g&%^$vvTUpKS%BM8z@*ru2#Bhc6x7&",
		"9@xKp!euwSwMWEImKlgIV@bBUVZas3yhRfu9rXtZsDW#oTmsu6cj6Tzoy7nGk0n9",
		"DUArkOShMP*UU@Nm1%W^FIz7kfawfknXlCfpcZ#9uqHgTfJT5$l4mres!b^dhdFF",
		"7qBiuE8afz&w@QCee1!gPTZEWe%6BJBEH1qCup90P$c%kSYuwivrk$c7#*Ms9z*D",
		"Jl0^$&Wt6flMz3d@iK0#RJeMSD4YS@Y!*^PpNH93fygT0$%zRQsdDbrXlXO5P4Au",
		"q2rUBie7uVp*@@Z2M%GN#vsWHPswSCt@^pGljvt9Qii6jU2u5*3jkQwioUjhFDCS",
		"bJRH7W3HjK8HlvYJF9JwqUCUR*CuHBIrbykZ7AUsNhj6SR5HDvxLdRc1VeHj#pjX",
		"*zvrQR#G*ngO!B0S3S#4dHasAYDODD9#8Qe%qb3NadsdPIJbFwy*oHQ7Rxakx@70",
		"eQLjyVFjq5z#7@dnTmPhtdHQ1eh5Rz3XA0AAdCfq!a0W5xlqo3pOLdht4xleGsmi",
		"A08vjv4kZqq#7YW2ZC#H!vrKEUTbzVwuIfNFJvkZG$pbh!W!aloVnfxRGhX*j71e"}
	c := &Crypto{
		secrctSeeds: seeds,
	}
	return c
}

//MakeMd5 用MD5加密后的字符串	DICK0 和 MAC用MD5够了
func (p *Crypto) makeMd5(rawstring string) string {
	md := md5.New()
	md.Write([]byte(rawstring))
	return hex.EncodeToString(md.Sum(nil))
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
	sh.Write([]byte(rawstring))
	return hex.EncodeToString(sh.Sum(nil))
}

//MakeSHA256 用SHA256加密后的字符串
func (p *Crypto) makeSHA256(rawstring string) string {
	sh := sha512.New512_256()
	sh.Write([]byte(rawstring))
	return hex.EncodeToString(sh.Sum(nil))
}

//MakePasword 密码加密 MD5 + SHA384 + BCRYPT
//这里只是展示客户端密码如何生成
func (p *Crypto) makePasword(password string) string {
	md := md5.New()
	sh := sha512.New() //只用到512
	md.Write([]byte(password))
	sh.Write(md.Sum(nil))
	bcryptstr, err := bcrypt.GenerateFromPassword(sh.Sum(nil), bcrypt.DefaultCost)
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
