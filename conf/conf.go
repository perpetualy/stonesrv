package conf

import (
	"gopkg.in/ini.v1"
	"stonesrv/models"
	"stonesrv/log"
	"fmt"
)

var conf = initConfig()

func GetDBAddress()string{
	return conf.DBAddress
}

func GetServerAddress()string{
	return conf.ServerAddress
}

func initConfig() models.Config{
	config,err := readConfig("conf/stonesrv.cfg")  //也可以通过os.arg或flag从命令行指定配置文件路径
	if err != nil {
		log.Error(fmt.Sprintf("%v",err))
		panic(-1)
	}
	log.Info(fmt.Sprintf("Config file read %+v",config))
	return config
}
//读取配置文件并转成结构体
func readConfig(path string) (models.Config, error) {
	var config models.Config
	conf, err := ini.Load(path)   //加载配置文件
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
