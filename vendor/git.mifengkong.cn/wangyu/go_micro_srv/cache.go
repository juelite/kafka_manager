package go_micro_srv_client

import (

	"golang.org/x/net/context"
	"fmt"
	"time"
)

//log在consul服务名称
const CACHE_SRV_NAME = "cache"

/**
 * 写入缓存
 */
func Cache(name string , express int64 , value string) error {

	request := &SetRequest{
		Key:name,
		Express:express,
		Val:value,
	}
	_, err := callCacheSrvSet(request)
	return err
}


/**
 * 获取缓存
 */
func GetCache(name string) (map[string]string , error) {

	request := &GetRequest{
		Key:name,
	}
	res, err := callCacheSrvGet(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	return res.Data, err
}
/**
 * 清除缓存
 */
func DelCache(name string) (error) {
	request := &DelRequest{
		Key:name,
	}
	_, err := callCacheSrvDel(request)
	return err
}


/**
 * 调用服务
 * @param request *FrLogRequest 请求入参
 * @return response *FrLogReply 返回结果
 * @return err error 返回错误信息
 */
func callCacheSrvSet(request *SetRequest) (response *SetResponse, err error) {
	//获取服务实例
	conn, err := GetSrvHandel(CACHE_SRV_NAME)
	if conn == nil || err != nil {
		return response, err
	}
	//初始化客户端
	c := NewCacheClient(conn)
	//请求服务
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	response, err = c.Set(ctx, request)
	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}
	return response, err
}

/**
 * 调用服务
 * @param request *FrLogRequest 请求入参
 * @return response *FrLogReply 返回结果
 * @return err error 返回错误信息
 */
func callCacheSrvGet(request *GetRequest) (response *GetResponse, err error) {
	//获取服务实例
	conn, err := GetSrvHandel(CACHE_SRV_NAME)
	if conn == nil || err != nil {
		return response, err
	}
	//初始化客户端
	c := NewCacheClient(conn)
	//请求服务
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	response, err = c.Get(ctx, request)
	if err != nil {
		return response, err
	}
	return response, err
}

/**
 * 调用服务
 * @param request *FrLogRequest 请求入参
 * @return response *FrLogReply 返回结果
 * @return err error 返回错误信息
 */
func callCacheSrvDel(request *DelRequest) (response *DelResponse, err error) {
	//获取服务实例
	conn, err := GetSrvHandel(CACHE_SRV_NAME)
	if conn == nil || err != nil {
		return response, err
	}
	//初始化客户端
	c := NewCacheClient(conn)
	//请求服务
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	response, err = c.Del(ctx, request)
	if err != nil {
		return response, err
	}
	return response, err
}