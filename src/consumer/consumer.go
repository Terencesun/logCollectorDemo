package consumer

import (
	collectorConfig "logCollector/src/config"
	collectorKafka "logCollector/src/kafka"
)

func InitConsumer(kafka collectorConfig.Kafka) (consumerInstance *collectorKafka.KafkaConsumer, err error) {
	consumerInstance, err = collectorKafka.InitKafkaConsumer(kafka)
	if err != nil {
		return
	}
	return
}
