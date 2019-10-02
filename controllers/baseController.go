package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"kafka_manager/constants"
)

type BaseController struct {
	beego.Controller
}

//返回结构声明
type ReturnData struct {
	code int
	message string
	data map[string]interface{}
}

/**
  * 接受参数方法
  * @param param interface{} 对应struct地址
 */
func (this *BaseController) RequestData(param interface{})  {
	this.ParseForm(param)
	// 这里的错误信息可以进行重写,学会了反射之后进行自定义错误信息
	valida := validation.Validation{}
	valida.Valid(param)
	message := ""
	if valida.HasErrors() {
		for _, value := range valida.Errors {
			message = value.Key + " " + value.Message
			break
		}
		var data map[string]interface{}
		this.ReturnData(constants.ERROR_CODE, message, data)
	}
}

/**
  * 返回json数据，请求结束
  * @code 状态码
  * @message 错误信息
  * data 数据
 */
func (this *BaseController) ReturnData(code int, message string, data interface{})  {
	this.Data["json"] =  map[string]interface{}{
		"code": code,
		"message" : message,
		"data" :data,
	}
	this.ServeJSON()
	this.StopRun()
}

