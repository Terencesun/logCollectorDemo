package main

import (
	collectorConfig "logCollector/src/config"
	collectorFactory "logCollector/src/factory"
	collectorLog "logCollector/src/log"
	"sync"
)


func main() {

	var waitGroup sync.WaitGroup

	conf, etcdIns, err := collectorConfig.InitConfig()
	if err != nil {
		panic("config init error")
	}

	err = collectorLog.InitLogger(conf)
	if err != nil {
		panic("logger init error")
	}

	factory := &collectorFactory.Factory{}
	factory.Init(conf, &waitGroup)

	for _, instance := range conf.Instances {
		waitGroup.Add(2)
		err := factory.AddWorker(instance, &waitGroup)
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

	waitGroup.Wait()
}
