package go_micro_srv_client

import (
	"github.com/0x5010/grpcp"
	"github.com/astaxie/beego"
	"google.golang.org/grpc"
	consul "github.com/hashicorp/consul/api"

	"fmt"
	"sync"
)


const (
	//consul 地址
	CONSUL_ADDR_PROD = "consul.mifengkong.cn:8500"
	CONSUL_ADDR_DEV = "demo.consul.mifengkong.cn"
	//consul token
	CONSUL_TOKEN_PROD = "70bfbc45-bb86-47a3-b786-9d16fc02b100"
	CONSUL_TOKEN_DEV = ""
)


var (
	env string
	pool *grpcp.ConnectionTracker
	service []*consul.CatalogService
	consul_addr, concul_token string
	client *consul.Client
)

func init() {
	so := sync.Once{}
	so.Do(func() {
		//初始化client
		var err error
		if getEnv() == "prod" {
			consul_addr = CONSUL_ADDR_PROD
			concul_token = CONSUL_TOKEN_PROD
		} else {
			consul_addr = CONSUL_ADDR_DEV
			concul_token = CONSUL_TOKEN_DEV
		}
		config := consul.Config{
			Address:consul_addr,
			Token:concul_token,
		}
		client, err = consul.NewClient(&config)//非默认情况下需要设置实际的参数
		if err != nil {
			client = nil
			fmt.Printf("did not connect: %v", err.Error())
		}

		//初始化连接池
		pool = grpcp.New(func(addr string) (*grpc.ClientConn, error) {
			return grpc.Dial(
				addr,
				grpc.WithInsecure(),
			)
		})
	})
}

/**
 * 获取环境变量
 */
func getEnv() string {
	env = beego.AppConfig.String("runmode")
	return env
}

/**
 * 获取服务
 * @param name string 服务名称
 * @return service 服务实例
 */
func GetSrvHandel(name string) (conn *grpc.ClientConn, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	service, _, err = client.Catalog().Service(name, "", nil)
	if err != nil {
		fmt.Printf("can not find services: %v", err.Error())
		return nil, err
	}
	//多服务端时遍历链接直到成功
	for _, v := range service  {
		addr := fmt.Sprintf("%s:%d", v.ServiceAddress, v.ServicePort)
		conn, err = pool.GetConn(addr)
		if err == nil {
			break
		}
	}
	return conn, err
}