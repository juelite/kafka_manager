package controllers

import (
	"github.com/astaxie/beego/orm"
	"kafka_manager/constants"
	"kafka_manager/models"
	"kafka_manager/services"
)

type ConsumeController struct {
	BaseController
}

type NewConsumeParam struct {
	TopicCallbackId int    `form:"topic_callback_id" valid:"Required"`
	ConsumerGroup   string `form:"consumer_group" valid:"Required;MinSize(3);MaxSize(64);Match(/^([a-zA-Z0-9-_]+)$/)"`
}

type DeleteConsumeParam struct {
	ConsumerId string `form:"consumer_id" valid:"Required"`
}

// 新增一个消费者
func (this *ConsumeController) NewConsumer() {
	param := &NewConsumeParam{}
	this.RequestData(param)
	topicCallModel := &models.TopicCallbackTemplate{}
	topicCallModel.Id = param.TopicCallbackId
	data, err := topicCallModel.GetTopicCallBackById()
	if err == orm.ErrNoRows {
		this.ReturnData(constants.ERROR_PARAMS, "不存在对应的模板id", nil)
	}
	consumerManagerService := &services.ConsumerManagerService{}
	err = consumerManagerService.StartConsumer(data.Topic, param.ConsumerGroup, data.CallbackUrl, true, "", param.TopicCallbackId)
	if err != nil {
		this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	this.ReturnData(constants.SUCCESS_CODE, "", nil)
}

// 获取当前所有消费者
func (this *ConsumeController) GetConsumers() {
	consumerManagerService := &services.ConsumerManagerService{}
	resp, err := consumerManagerService.GetAllConsumersData()
	if err != nil {
		this.ReturnData(constants.SERVER_ERROR, constants.SERVER_ERROR_MSG, nil)
	}
	var res = make([]*services.ManagerConsumer, 0, 100)
	for _, consumeRecord := range resp {
		res = append(res, consumeRecord)
	}
	this.ReturnData(constants.SUCCESS_CODE, "", res)
}

// 删除某个消费者
func (this *ConsumeController) DeleteConsumer() {
	param := &DeleteConsumeParam{}
	this.RequestData(param)
	consumerManagerService := &services.ConsumerManagerService{}
	if ok, err := consumerManagerService.StopOneConsume(param.ConsumerId); ok == false {
		this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	this.ReturnData(constants.SUCCESS_CODE, "", nil)
}
