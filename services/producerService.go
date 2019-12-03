package services

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka_manager/common"
	"kafka_manager/common/kafka_aliware"
)

type ProduceService struct{}

var producer sarama.SyncProducer

// 生产者初始化
func init() {
	var err error
	kafkaConfig := kafka_aliware.KafkaConfig{}
	producer, err = sarama.NewSyncProducer(kafkaConfig.GetServers(), kafkaConfig.GetProducerConfig())
	fmt.Println("生产者开始运行，调用/api/message/new接口生产信息")
	if err != nil {
		panic("Kafka生产者初始化失败" + err.Error())
	}
}

// 生产信息到消息队列
func (this *ProduceService) Produce(topic string, traceId string, value string) (err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(traceId),
		Value: sarama.StringEncoder(value),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		// 记入kafka失败，存入日志
		common.TraceLogger(traceId, "生产者写入kafka失败"+value, "error")
		return
	}
	return
}
