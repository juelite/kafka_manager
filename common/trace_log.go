package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"git.mifengkong.cn/wangyu/go_micro_srv"
	"github.com/astaxie/beego"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// 链路追踪，跟踪信息的整个流程
func TraceLogger(traceId string, value string, level string) {
	// 微服务日志记录
	go_micro_srv_client.ElcLog(beego.AppConfig.String("kafkaTraceKey"), traceId+" -> trace value : "+value, level)
}

// 生成traceId
func GenerateTraceId(key string) (value string) {
	timeSalt := time.Now().Unix()
	key = key + fmt.Sprint(timeSalt) + fmt.Sprint(rand.Intn(10000))
	encodeObject := sha1.New()
	encodeObject.Write([]byte(key))
	return hex.EncodeToString(encodeObject.Sum([]byte("")))
}
