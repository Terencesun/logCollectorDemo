package config

import (
	"encoding/json"
	"fmt"
	beegoConfig "github.com/astaxie/beego/config"
	"path/filepath"
)

type Instance struct {
	Topic       string `json:"topic"`
	LogFilePath string `json:"logPath"`
}

type Kafka struct {
	Hosts []string
}

type ConfigStore struct {
	SystemLogPath  string     // 系统日志地址
	SystemLogLevel string     // 系统日志等级
	Instances      []Instance // 监听实例地址
	KafkaConfig    Kafka      // kafka配置
}

func InitConfig() (conf ConfigStore, err error) {
	conf = ConfigStore{}
	// 初始化地址
	confPath, err := filepath.Abs("./env.conf")
	fmt.Println(confPath)
	if err != nil {
		err = CONFPATHERROR
		return
	}

	// 初始化config
	configer, err := beegoConfig.NewConfig("ini", confPath)
	if err != nil {
		err = CONFINITERROR
		return
	}

	serverConf := configer.String("server::instances")
	err = json.Unmarshal([]byte(serverConf), &(conf.Instances))
	if err != nil {
		err = CONFINITERROR
		return
	}

	conf.SystemLogPath, _ = filepath.Abs(configer.String("server::serverLogPath"))
	conf.SystemLogLevel, _ = filepath.Abs(configer.String("server::serverLogLevel"))

	kafkaConf := configer.String("kafka::hosts")
	err = json.Unmarshal([]byte(kafkaConf), &(conf.KafkaConfig.Hosts))
	if err != nil {
		err = CONFINITERROR
		return
	}

	return
}
