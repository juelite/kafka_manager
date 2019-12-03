package controllers

import (
	"encoding/json"
	"kafka_manager/common"
	"kafka_manager/services"
)

type ProduceController struct {
	BaseController
}

// 生产信息，写入kafka消息队列，主要功能为削峰填谷
func (this *ProduceController) ProduceMessage() {
	topic := this.GetString("topic")
	if topic == "" {
		this.ReturnData(1, "topic参数缺少", nil)
	}
	method := this.Ctx.Request.Method
	if method != "GET" && method != "POST" {
		this.ReturnData(1, "请求方法不支持，仅支持get和post", nil)
	}
	params := this.Input()
	messageBody := services.MessageQueueData{}
	messageBody.RequestMethod = method
	messageBody.Params = messageBody.BuildQueryString(params)
	value, err := json.Marshal(messageBody)
	if err != nil {
		this.ReturnData(1, "参数错误或者参数格式错误", nil)
	}
	// 将该有的日志信息全部记录下来, ip：端口号 请求时间 所有参数 请求的域名
	ip := this.Ctx.Request.RemoteAddr
	host := this.Ctx.Input.Host()

	// 根据主题生成trace id,并初始化链路
	traceId := common.GenerateTraceId(topic)
	common.TraceLogger(traceId, "消息请求初始化(请求到达) "+"IP: "+ip+":"+" Hosts: "+host+" method:"+method+" params: "+messageBody.Params, "info")

	produceService := &services.ProduceService{}
	err = produceService.Produce(topic, traceId, string(value))
	if err != nil {
		this.ReturnData(1, "消息队列异常，请稍后再试", nil)
	}
	common.TraceLogger(traceId, "消息请求结束(消息成功写入队列) "+"IP: "+ip+" Hosts: "+host+" method:"+method+" params: "+messageBody.Params, "info")
	// 将追踪id返回给请求方，可以追踪整个消费场景
	this.ReturnData(0, "请求成功", traceId)
}
