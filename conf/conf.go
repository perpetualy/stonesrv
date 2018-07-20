package conf

import (
	"fmt"
	"stonesrv/log"
	"stonesrv/models"
	"strings"

	"gopkg.in/ini.v1"
)

var conf = newConfig()

//newConfig 新建配置文件对象
func newConfig() *Conf {
	c := &Conf{
		confPath: "./conf/stonesrv.cfg",
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

//GetServerAddress 获取服务器地址
func GetServerAddress() string {
	return conf.config.ServerAddress
}

//GetServerPort 获取服务器端口
func GetServerPort() string {
	return conf.config.ServerPort
}

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
