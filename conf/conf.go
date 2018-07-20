package conf

import (
	"gopkg.in/ini.v1"
	"stonesrv/models"
	"stonesrv/log"
	"fmt"
)

var conf = newConfig()

func newConfig() *Conf{
	c := &Conf{
		confPath : "conf/stonesrv.cfg",
	}
	c.config = c.initConfig()
	return c
}

type Conf struct{
	confPath string
	config   models.Config
}

func GetDBAddress() string{
	return conf.config.DBAddress
}

func GetServerAddress() string{
	return conf.config.ServerAddress
}

func GetServerPort() string{
	return conf.config.ServerPort
}

//初始化配置文件
func (p *Conf)initConfig() models.Config{
	config,err := p.readConfig() 
	if err != nil {
		log.Error(fmt.Sprintf("%v",err))
		panic(-1)
	}
	log.Info(fmt.Sprintf("Config file read %+v",config))
	return config
}

//读取配置文件并转成结构体
func (p *Conf)readConfig() (models.Config, error) {
	var config models.Config
	conf, err := ini.Load(p.confPath)   //加载配置文件
	if err != nil {
		log.Error("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)   //解析成结构体
	if err != nil {
		log.Error("mapto config file fail!")
		return config, err
	}
	return config, nil
}
