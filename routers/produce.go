package routers

import (
	"github.com/astaxie/beego"
	"kafka_manager/controllers"
)

func init() {
	// 生产者生产消息(对外)
	beego.Router("/message/new", &controllers.ProduceController{}, "*:ProduceMessage")
}
