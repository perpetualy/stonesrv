package conf

import (
	"fmt"
	"strings"
	"xmvideo/log"
	"xmvideo/models"

	"gopkg.in/ini.v1"
)

var conf = newConfig()

//newConfig 新建配置文件对象
func newConfig() *Conf {
	c := &Conf{
		confPath: "./conf/xmvideosrv.cfg",
	}
	return c
}

//Conf 配置文件
type Conf struct {
	confPath string
	config   models.Config
}

//Init 初始化方法
func Init(path string) {
	conf = newConfig()
	if strings.Compare(path, "") != 0 {
		conf.confPath = path
	}
	conf.initConfig()
}

//GetDBAddress 获取数据库地址
func GetDBAddress() string {
	return conf.config.DBAddress
}

//GetDBUser 获取配置文件数据库登陆用户名
func GetDBUser() string {
	//env.ArangoDBDefaultUser
	return conf.config.DBUser
}

//GetDBPassword 获取配置文件数据库登陆密码
func GetDBPassword() string {
	//env.ArangoDBDefaultPassword
	return conf.config.DBPassword
}

//GetDBName 获取配置文件数据库登陆密码
func GetDBName() string {
	//env.ArangoDBDefaultDBName
	return conf.config.DBName
}

//GetServerAddress 获取服务器地址
func GetServerAddress() string {
	return conf.config.ServerAddress
}

//GetServerPort 获取服务器端口
func GetServerPort() string {
	return conf.config.ServerPort
}

//GetSSLCrtFile 获取SSL CRT证书
func GetSSLCrtFile() string {
	return conf.config.SSLCrtFile
}

//GetSSLKeyFile 获取SSL KEY
func GetSSLKeyFile() string {
	return conf.config.SSLKeyFile
}

//GetUpdatesDir 获取更新路径
func GetUpdatesDir() string {
	return conf.config.UpdatesDir
}

//GetUpdateFile 获取更新文件名
func GetUpdateFile() string {
	return conf.config.UpdateFile
}

//GetLanguage 获取语言
func GetLanguage() string {
	return conf.config.Language
}

//IsDebugMode 判断是否DEBUG模式
func IsDebugMode() bool {
	return strings.Compare(conf.config.DebugMode, "1") == 0
}

//IsToCMode 判断是否ToC模式
// func IsToCMode() bool {
// 	return strings.Compare(conf.config.ToCMode, "1") == 0
// }

//initConfig 初始化配置文件
func (p *Conf) initConfig() {
	config, err := p.readConfig()
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		panic(-1)
	}
	log.Info(fmt.Sprintf("Config file read %+v", config))
	p.config = config
}

//readConfig 读取配置文件并转成结构体
func (p *Conf) readConfig() (models.Config, error) {
	var config models.Config
	conf, err := ini.Load(p.confPath) //加载配置文件
	if err != nil {
		log.Error("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		log.Error("mapto config file fail!")
		return config, err
	}
	return config, nil
}
