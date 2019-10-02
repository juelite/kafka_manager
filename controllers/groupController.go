package controllers

import (
	"github.com/Shopify/sarama"
	"kafka_manager/common/kafka_aliware"
	"kafka_manager/constants"
)

type GroupController struct {
	BaseController
}

// 获取集群的consumer groups
func (this *GroupController) GetConsumeGroups() {
	kafkaConf := &kafka_aliware.KafkaConfig{}
	// 由于初调查，sarama不支持集群查询consumeGroup，无法确定哪个broker是controller，故轮询
	var data = make([]map[string]string, 0, 100)
	for _, server := range kafkaConf.GetServers() {
		client := sarama.NewBroker(server)
		err := client.Open(kafkaConf.GetConsumeGroupConfig())
		if err != nil {
			this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
		}
		rsp, err := client.ListGroups(&sarama.ListGroupsRequest{})
		if err != nil {
			this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
		}
		for groupName, groupType:= range rsp.Groups{
			var tmp = make(map[string]string)
			tmp["name"] = groupName
			tmp["type"] = groupType
			data = append(data, tmp)
		}
		client.Close()
	}

	this.ReturnData(constants.SUCCESS_CODE, "", data)
}

// sarama api不支持，暂废弃
func (this *GroupController) CreateConsumeGroups()  {

}
