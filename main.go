package main

import (
	collectorConfig "logCollector/src/config"
	collectorConsumer "logCollector/src/consumer"
	collectorEtcd "logCollector/src/etcd"
	collectorLog "logCollector/src/log"
	collectorProducer "logCollector/src/producer"
	"sync"
)

var (
	systemEnv map[string]string
)

func main() {

	systemEnv = make(map[string]string)
	var waitGroup sync.WaitGroup

	conf, err := collectorConfig.InitConfig()
	if err != nil {
		panic("config init error")
	}

	err = collectorLog.InitLogger(conf)
	if err != nil {
		panic("logger init error")
	}

	err = collectorEtcd.InitEtcd(conf)
	if err != nil {
		panic("logger init error")
	}

	for _, instance := range conf.Instances {
		waitGroup.Add(2)
		tailInstance, kafkaProInstance, err := collectorProducer.InitTailAndKafka(instance, conf.KafkaConfig)
		kafkaConInstance, err := collectorConsumer.InitConsumer(conf.KafkaConfig)
		if err != nil {
			panic("tail init error")
		}
		systemEnv[instance.LogFilePath] = instance.Topic
		tailInstance.Start(kafkaProInstance, kafkaConInstance, &waitGroup)
	}

	waitGroup.Wait()
}
