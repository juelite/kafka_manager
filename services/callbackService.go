package services

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type CallbackService struct {
}

type MessageQueueData struct {
	RequestMethod string
	Params        string
}

// 回调的处理方法
func (this *CallbackService) ProcessMessage(callback string, key string, value string) (err error) {
	if !strings.Contains(callback, "http") {
		callback = "http://" + callback
	}
	valueByte := []byte(value)
	messageData := &MessageQueueData{}
	json.Unmarshal(valueByte, messageData)
	var req *httplib.BeegoHTTPRequest
	if messageData.RequestMethod == "GET" {
		req = httplib.Get(callback)
	} else if messageData.RequestMethod == "POST" {
		req = httplib.Post(callback)
	} else {
		return errors.New("请求不能为GET和POST以外的")
	}
	paramsString := messageData.Params
	paramArr := strings.Split(paramsString, "&")
	for _, v := range paramArr {
		vArr := strings.SplitN(v, "=", 2) // todo
		if len(vArr) != 2 {
			return errors.New("数据源错误" + paramsString)
		}
		req.Param(vArr[0], vArr[1])
	}
	req.SetTimeout(time.Second*5, time.Second*5)
	_, err = req.Bytes()
	if err != nil {
		return err
	}
	return nil
}

// 生成url参数字符串
func (this *MessageQueueData) BuildQueryString(params map[string][]string) string {
	buf := new(bytes.Buffer)
	for k, v := range params {
		for _, v2 := range v {
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(v2)
			buf.WriteString("&")
		}
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}
