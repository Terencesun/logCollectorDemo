package main

import (
	"fmt"
	collectorApi "logCollector/src/api"
	collectorConfig "logCollector/src/config"
	collectorFactory "logCollector/src/factory"
	collectorLog "logCollector/src/log"
)


func main() {

	conf, etcdIns, err := collectorConfig.InitConfig()
	if err != nil {
		panic("config init error")
	}
	fmt.Println("config inited.")

	err = collectorLog.InitLogger(conf)
	if err != nil {
		panic("logger init error")
	}
	fmt.Println("logger inited.")

	factory := &collectorFactory.Factory{}
	factory.Init(conf)

	for _, instance := range conf.Instances {
		err := factory.AddWorker(instance)
		if err != nil {
			panic("worker init error")
		}
	}

	wch, err := etcdIns.SetWatch("server/instances")
	if err != nil {
		panic("watch init error")
	}

	go func(watchChan *chan string) {
		for v := range *watchChan {
			err := factory.UpdateWorker(v)
			if err != nil {
				panic(err)
			}
		}
	}(wch)

	go func() {
		err := collectorApi.Start("127.0.0.1:1234", &conf.EtcdConfig)
		if err != nil {
			panic("api init error")
		}
	}()

	var hold chan bool
	for {
		select {
		case <- hold:
			return
		}
	}
}
