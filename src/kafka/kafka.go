package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	collectorConfig "logCollector/src/config"
	"sync"
)

type KafkaProducer struct {
	Client sarama.SyncProducer
	KillChan *chan bool
}

func (p *KafkaProducer) SendMsg(topic string, value []byte) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(value)

	_, _, err = p.Client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed", err)
		return
	}
	return
}

func InitKafkaProducer(conf collectorConfig.Kafka, killChan *chan bool) (kafkaInstance *KafkaProducer, err error) {
	kafkaInstance = &KafkaProducer{
		KillChan: killChan,
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	kafkaInstance.Client, err = sarama.NewSyncProducer(conf.Hosts, config)
	if err != nil {
		return
	}
	go func() {
		loop: for {
			select {
			case <- *killChan:
				kafkaInstance.Client.Close()
				fmt.Printf("close kafka producer, %v\n", conf.Hosts)
				break loop
			}
		}
	}()
	return
}

type KafkaConsumer struct {
	Client        sarama.Consumer
	PartitionList []int32
}

func (p *KafkaConsumer) Init(topic string, killChan *chan bool) (err error) {
	var wg sync.WaitGroup
	partitionList, err := p.Client.Partitions(topic)
	if err != nil {
		return
	}
	defer p.Client.Close()
	p.PartitionList = partitionList
	for partition := range partitionList {
		partitionConsumer, pcErr := p.Client.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if pcErr != nil {
			return pcErr
		}
		defer partitionConsumer.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			loop: for {
				select {
				case msg := <-pc.Messages():
					// todo 写入es
					fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				case <- *killChan:
					fmt.Printf("close consumer, %v\n", topic)
					break loop
				}
			}
			wg.Done()
		}(partitionConsumer)
	}
	wg.Wait()
	return
}

func InitKafkaConsumer(conf collectorConfig.Kafka) (kafkaInstance *KafkaConsumer, err error) {
	kafkaInstance = &KafkaConsumer{}
	kafkaInstance.Client, err = sarama.NewConsumer(conf.Hosts, nil)
	if err != nil {
		return
	}
	return
}
