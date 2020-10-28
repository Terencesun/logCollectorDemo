package main

import (
	collectorConfig "./config"
	collectorConsumer "./consumer"
	collectorLog "./log"
	collectorTail "./tail"
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

	for _, instance := range conf.Instances {
		waitGroup.Add(2)
		tailInstance, kafkaProInstance, err := collectorTail.InitTailAndKafka(instance, conf.KafkaConfig)
		kafkaConInstance, err := collectorConsumer.InitConsumer(instance, conf.KafkaConfig)
		if err != nil {
			panic("tail init error")
		}
		systemEnv[instance.LogFilePath] = instance.Topic
		tailInstance.Start(kafkaProInstance, kafkaConInstance, &waitGroup)
	}

	waitGroup.Wait()
}
