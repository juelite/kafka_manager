package go_micro_srv_client

import (
	"golang.org/x/net/context"
	"fmt"
	"time"
)

//log在consul服务名称
const LOG_SRV_NAME = "kibana"

/**
 * 写elc日志
 * @param tag string 标签
 * @param info string 内容
 * @param level string 级别 默认info
 */
func ElcLog(tag string , info string , level string) error {
	var err error
	if level == "" {
		level = "info"
	}
	request := &WriteRequest{
		Tag:tag,
		Info:info,
		Level:level,
	}
	_, err = callLogSrv(request)
	return err
}

/**
 * 调用服务
 * @param request *FrLogRequest 请求入参
 * @return response *FrLogReply 返回结果
 * @return err error 返回错误信息
 */
func callLogSrv(request *WriteRequest) (response *WriteResponse, err error) {
	//获取服务实例
	conn, err := GetSrvHandel(LOG_SRV_NAME)
	if conn == nil || err != nil {
		return response, err
	}
	//初始化客户端
	c := NewKibanaClient(conn)
	//请求服务
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	response, err = c.Write(ctx, request)
	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}
	return response, err
}