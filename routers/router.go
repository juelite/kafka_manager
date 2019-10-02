package routers

import (
	"kafka_manager/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// vue单页面应用
    beego.Router("/", &controllers.MainController{})
}
