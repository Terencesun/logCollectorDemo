package factory

import (
	"encoding/json"
	"fmt"
	collectorConfig "logCollector/src/config"
	collectorConsumer "logCollector/src/consumer"
	collectorProducer "logCollector/src/producer"
	"sync"
)

type Factory struct {
	WorkerInfo map[string]*Worker
	Config collectorConfig.ConfigStore
	Wg *sync.WaitGroup
}

func (p *Factory) Init(conf collectorConfig.ConfigStore)  {
	p.WorkerInfo = make(map[string]*Worker)
	p.Config = conf
}

func (p *Factory) AddWorker(instance collectorConfig.Instance) (err error) {
	wk := &Worker{
		path: instance.LogFilePath,
	}

	p.WorkerInfo[instance.LogFilePath] = wk

	err = wk.Create(instance, p.Config)

	return
}

func (p *Factory) RemoveWorker(path string) (err error) {
	err = p.WorkerInfo[path].Destory()
	if err != nil {
		return
	}
	delete(p.WorkerInfo, path)
	return
}

func (p *Factory) UpdateWorker(instanceString string) (err error) {
	var instances []collectorConfig.Instance
	err = json.Unmarshal([]byte(instanceString), &instances)
	// 判断有没有需要删除的
	for key := range p.WorkerInfo {
		var isExist bool = false
		for _, v := range instances {
			if v.LogFilePath == key {
				isExist = true
				break
			}
		}
		if !isExist {
			path := p.WorkerInfo[key].path
			err = p.RemoveWorker(path)
			if err != nil {
				fmt.Printf("删除旧的实例%v失败\n", path)
				return
			}
			fmt.Printf("删除旧的实例%v\n", path)
		}
	}
	// 判断有没有新的，如果有新的，创建，如果不是新的，pass
	for _, v := range instances {
		if _, ok := p.WorkerInfo[v.LogFilePath]; !ok {
			// 新的，进行创建
			err = p.AddWorker(v)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("创建新的实例%v失败\n", v.LogFilePath)
				return
			}
			fmt.Printf("创建新的实例%v\n", v.LogFilePath)
		}
	}
	return
}

type Worker struct {
	path string
	killChan chan bool
}

func (p *Worker) Create(instance collectorConfig.Instance, conf collectorConfig.ConfigStore) (err error) {
	p.killChan = make(chan bool)
	tailInstance, kafkaProInstance, err := collectorProducer.InitTailAndKafka(instance, conf.KafkaConfig, &p.killChan)
	kafkaConInstance, err := collectorConsumer.InitConsumer(conf.KafkaConfig)
	if err != nil {
		err = TAILINITERROR
		return
	}
	tailInstance.Start(kafkaProInstance, kafkaConInstance)
	return
}

func (p *Worker) Destory() (err error) {
	close(p.killChan)
	return
}
