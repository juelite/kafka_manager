package routers

import (
	"github.com/astaxie/beego"
	"kafka_manager/controllers"
)

// kafka管理界面ajax 请求的接口
func init() {
	// 获取kafka集群所有的topics,可筛选
	beego.Router("/api/getTopics", &controllers.TopicController{}, "GET:GetTopicsData")
	// 获取某一个topic的分区状态
	//beego.Router("/api/getPartitions", &controllers.MainController{})
	// 创建一个topic
	beego.Router("/api/createTopic", &controllers.TopicController{}, "POST:CreateTopic")

	// 获取所有的consume group管理,可筛选
	beego.Router("/api/getConsumerGroups", &controllers.GroupController{}, "GET:GetConsumeGroups")
	// 查看某个consume group的消费堆积
	beego.Router("api/getConsumeGroupInfo", &controllers.MainController{})

	// 查询所有消费者模板
	beego.Router("/api/getConsumerTemplates", &controllers.TopicCallbackController{}, "GET:GetTopicCallbackTemplates")
	// 创建一个消费者模板
	beego.Router("/api/createConsumerTemplate", &controllers.TopicCallbackController{}, "POST:CreateTopicCallbackTemplate")
	// 新开启一个消费实例
	beego.Router("/api/startConsumer", &controllers.ConsumeController{}, "POST:NewConsumer")
	// 关闭某个消费实例
	beego.Router("/api/closeConsumer", &controllers.ConsumeController{}, "POST:DeleteConsumer")
	// 获取当前所有消费实例,可筛选
	beego.Router("/api/getConsumers", &controllers.ConsumeController{}, "GET:GetConsumers")

	// 获取模板下的所有消费者
	beego.Router("/api/getAllConsumersOfTemplate", &controllers.TopicCallbackController{}, "GET:GetAllConsumersOfTemplate")
	// 删除指定模板
	beego.Router("/api/deleteTemplate", &controllers.TopicCallbackController{}, "GET:DeleteTemplate")
}
