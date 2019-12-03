package controllers

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"kafka_manager/common/kafka_aliware"
	"kafka_manager/constants"
	"log"
)

type TopicController struct {
	BaseController
}

type CreateTopicParam struct {
	//1. Topic 名称只能包含字母，数字，下划线(_)和短横线(-)
	//2. 名称长度限制在 3-64 字节之间，长于 64 字节将被自动截取
	//3. 一旦创建后不能再修改 Topic 名称
	Topic string `form:"topic_name" valid:"Required;MinSize(3);MaxSize(64);Match(/^([a-zA-Z0-9-_]+)$/)"`
}

// 获取当前集群的所有topic
func (t *TopicController) GetTopicsData() {
	kafkaConf := &kafka_aliware.KafkaConfig{}
	client, err := sarama.NewClient(kafkaConf.GetServers(), kafkaConf.GetConfig())
	if err != nil {
		t.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
		return
	}
	defer client.Close()
	topics, err := client.Topics()
	if err != nil {
		t.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
		return
	}
	// 一个切片，然后切片的每一个元素为一个map,cap为100，避免频繁扩容
	var data = make([]map[string]string, 0, 100)
	for _, topicName := range topics {
		var tmp = make(map[string]string)
		tmp["name"] = topicName
		data = append(data, tmp)
	}
	t.ReturnData(constants.SUCCESS_CODE, "", data)
}

// 创建一个新的topic
func (t *TopicController) CreateTopic() {
	param := &CreateTopicParam{}
	t.RequestData(param)
	kafkaConf := kafka_aliware.KafkaConfig{}
	admin, err := sarama.NewClusterAdmin(kafkaConf.GetServers(), kafkaConf.GetAdminClusterConfig())
	if err != nil {
		// 将错误日志打出来,可以使用它们的日志微服务
		t.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
		return
	}
	defer admin.Close()
	// 这个后期做成可自定义配置
	numPartitions, _ := beego.AppConfig.Int("TOPIC_PARTITIONS")
	replicationFactor, _ := beego.AppConfig.Int("REPLICATION_FACTOR")
	topicDetail := sarama.TopicDetail{
		NumPartitions:     int32(numPartitions),
		ReplicationFactor: int16(replicationFactor),
	}
	// 集群创建topic，不区分单台机器
	err = admin.CreateTopic(param.Topic, &topicDetail, false)
	if err != nil {
		log.Println(2)
		t.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	t.ReturnData(constants.SUCCESS_CODE, "", nil)
}
