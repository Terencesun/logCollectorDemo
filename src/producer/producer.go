package producer

import (
	"fmt"
	"github.com/hpcloud/tail"
	collectorConfig "logCollector/src/config"
	collectorKafka "logCollector/src/kafka"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TailInfo struct {
	LogPath string
	Topic   string
	tails   *tail.Tail
}

func (p *TailInfo) Init() (err error) {
	filename := p.LogPath
	p.tails, err = tail.TailFile(filename, tail.Config{
		ReOpen:    true,                                           // 是否重新打开
		Follow:    true,                                           // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}, // 文件从哪个地方开始读
		MustExist: false,                                          // 文件不存在不报错
		Poll:      true,                                           // 文件变化查询的方式，poll是轮询，还有inotify，通知用户文件变化
	})

	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}

	return
}

func (p *TailInfo) Start(kafkaProInstance *collectorKafka.KafkaProducer, kafkaConInstance *collectorKafka.KafkaConsumer, lock *sync.WaitGroup) {
	go func() {
		var msg *tail.Line
		var ok bool
		for {
			msg, ok = <-p.tails.Lines
			if !ok {
				fmt.Printf("tail file close reopen, filename:%s\n", p.tails.Filename)
				time.Sleep(time.Second * 1)
				continue
			}
			// 写入kafka
			fmt.Println("msg:", msg.Text)
			err := kafkaProInstance.SendMsg(p.Topic, []byte(msg.Text))
			if err != nil {
				continue
			}
		}
		lock.Done()
	}()
	go func() {
		err := kafkaConInstance.Init(p.Topic)
		if err != nil {
			lock.Done()
		}
		lock.Done()
	}()
}

func (p *TailInfo) Kill() {

}

func InitTailAndKafka(instance collectorConfig.Instance, kafka collectorConfig.Kafka) (tailInstance *TailInfo, kafkaProInstance *collectorKafka.KafkaProducer, err error) {
	logPath, err := filepath.Abs(instance.LogFilePath)
	tailInstance = &TailInfo{
		LogPath: logPath,
		Topic:   instance.Topic,
	}

	kafkaProInstance, err = collectorKafka.InitKafkaProducer(kafka)

	err = tailInstance.Init()

	return
}
