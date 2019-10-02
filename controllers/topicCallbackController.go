package controllers

import (
	"kafka_manager/constants"
	"kafka_manager/models"
)

type TopicCallbackController struct {
	BaseController
}

type CreateTopicCallbackParam struct {
	Topic       string `form:"topic_name" valid:"Required;MinSize(3);MaxSize(64);Match(/^([a-zA-Z0-9-_]+)$/)"`
	CallbackUrl string `form:"callback_url" valid:"Required;"`
}

type GetConsumersOfTemplateParam struct {
	TemplateId int `form:"template_id" valid:"Required"`
}

type DeleteTemplateParam struct {
	TemplateId int `form:"template_id" valid:"Required"`
}

// 批量查询消费者模板
func (this *TopicCallbackController) GetTopicCallbackTemplates() {
	model := models.TopicCallbackTemplate{}
	res, err := model.GetAllTopicCallbacks()
	if err != nil {
		this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	data := make([]interface{}, 0, 100)
	for _, v := range res {
		data = append(data, v)
	}
	this.ReturnData(constants.SUCCESS_CODE, "", data)
}

// 创建一个消费者模板
func (this *TopicCallbackController) CreateTopicCallbackTemplate() {
	param := &CreateTopicCallbackParam{}
	this.RequestData(param)
	templateData := models.TopicCallbackTemplate{}
	templateData.Topic = param.Topic
	templateData.CallbackUrl = param.CallbackUrl
	ok, err := templateData.IsExistTopicAndCallback()
	if err != nil {
		this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	// 当前模板已经存在
	if ok == true {
		this.ReturnData(constants.TOPIC_CALLBACK_HAS_EXIST, constants.TOPIC_CALLBACK_HAS_EXIST_MSG, nil)
	}
	_, err = templateData.Insert()
	if err != nil {
		this.ReturnData(constants.SERVER_ERROR, err.Error(), nil)
	}
	this.ReturnData(constants.SUCCESS_CODE, "", nil)
}

// 获取模板开启的所有消费者
func (this *TopicCallbackController) GetAllConsumersOfTemplate() {
	param := &GetConsumersOfTemplateParam{}
	this.RequestData(param)
	backUpModel := &models.ConsumerBackup{}
	data := backUpModel.GetAllLiveConsumerInfoByTemplateId(param.TemplateId)
	this.ReturnData(0, "success", data)
}

// 删除指定的模板
func (this *TopicCallbackController) DeleteTemplate() {
	param := &GetConsumersOfTemplateParam{}
	this.RequestData(param)
	// 判断是否全部删除
	backUpModel := &models.ConsumerBackup{}
	if backUpModel.IsConsumerShutdownOfTemplateId(param.TemplateId) == false {
		this.ReturnData(1, "该模板还有其余消费者没有关闭，不能删除", nil)
	}
	model := &models.TopicCallbackTemplate{}
	model.DeleteByTemplateId(param.TemplateId)
	this.ReturnData(0, "success", nil)
}
