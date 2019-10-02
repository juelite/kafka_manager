package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_ "kafka_manager/boots"
	_ "kafka_manager/constants"
	_ "kafka_manager/routers"
	"kafka_manager/services"
)

func main() {
	// 本地测试和前端联调存在跨域问题
	if beego.AppConfig.String("runmode") != "prod" {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,
		}))
	}

	// 重启后还原消费者现场
	consumerLive := &services.ConsumerManagerService{}
	consumerLive.InitConsumerAndRestore()

	beego.Run()
}
