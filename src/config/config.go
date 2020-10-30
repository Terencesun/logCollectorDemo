package config

import (
	"encoding/json"
	"fmt"
	beegoConfig "github.com/astaxie/beego/config"
	collectorEtcd "logCollector/src/etcd"
	"path/filepath"
)

type Instance struct {
	Topic       string `json:"topic"`
	LogFilePath string `json:"logPath"`
}

type Kafka struct {
	Hosts []string
}

type Etcd struct {
	Hosts []string
}

type ConfigStore struct {
	SystemLogPath  string     // 系统日志地址
	SystemLogLevel string     // 系统日志等级
	Instances      []Instance // 监听实例地址
	KafkaConfig    Kafka      // kafka配置
	EtcdConfig	   Etcd		  // etcd配置
}

func InitConfig() (conf ConfigStore, etcdIns *collectorEtcd.Etcd, err error) {
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

	etcdConf := configer.String("etcd::hosts")
	err = json.Unmarshal([]byte(etcdConf), &(conf.EtcdConfig.Hosts))
	if err != nil {
		err = CONFINITERROR
		return
	}

	etcdIns, err = collectorEtcd.InitEtcd(conf.EtcdConfig.Hosts)
	if err != nil {
		panic("logger init error")
	}

	serverLogPath, err := etcdIns.GetKey("server/serverLogPath")
	conf.SystemLogPath, _ = filepath.Abs(serverLogPath)

	serverLevel, err := etcdIns.GetKey("server/logLevel")
	conf.SystemLogLevel, _ = filepath.Abs(serverLevel)

	serverInses, err := etcdIns.GetKey("server/instances")
	err = json.Unmarshal([]byte(serverInses), &(conf.Instances))
	if err != nil {
		err = CONFINITERROR
		return
	}

	kafkaConf, err := etcdIns.GetKey("kafka/hosts")
	err = json.Unmarshal([]byte(kafkaConf), &(conf.KafkaConfig.Hosts))
	if err != nil {
		err = CONFINITERROR
		return
	}

	return
}
