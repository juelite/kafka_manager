package kafka_aliware

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/bsm/sarama-cluster"
	"io/ioutil"
	"strings"
	"time"
)

var MqConfig *sarama.Config
var MqClusterConfig *cluster.Config

type KafkaConfig struct {

}

// saramaConfig公共的配置部分
func saramaConfigCommon() *sarama.Config{
	MqConfig = sarama.NewConfig()
	// 测试环境进行sasl认证，阿里云kafka要求
	if beego.AppConfig.String("runmode") == "test" {
		MqConfig.Net.SASL.Enable = true
		MqConfig.Net.SASL.User = beego.AppConfig.String("USERNAME")
		MqConfig.Net.SASL.Password = beego.AppConfig.String("PASSWORD")
		MqConfig.Net.SASL.Handshake = true

		MqConfig.Net.TLS.Enable = true
		certBytes, err := ioutil.ReadFile("./conf/ca-cert")
		if err != nil {
			panic("kafka manager failed fo find ca-cert.file, please check directory")
		}
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			panic("kafka manager config failed to parse root certificate")
		}
		MqConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            clientCertPool,
			InsecureSkipVerify: true,
		}
	}
	// 生产环境已确保Linux网络环境调用安全，所以不需要SASL校验
	if beego.AppConfig.String("runmode") == "prod" {
		MqConfig.Net.SASL.Enable = false
		MqConfig.Net.TLS.Enable = false
	}
	return MqConfig
}

// clusterConfig公共的配置部分
func clusterConfigCommon() *cluster.Config {
	MqClusterConfig = cluster.NewConfig()
	// 测试环境进行sasl认证，阿里云kafka要求
	if beego.AppConfig.String("runmode") == "test" {
		MqClusterConfig.Net.SASL.Enable = true
		MqClusterConfig.Net.SASL.User = beego.AppConfig.String("USERNAME")
		MqClusterConfig.Net.SASL.Password = beego.AppConfig.String("PASSWORD")
		MqClusterConfig.Net.SASL.Handshake = true

		MqClusterConfig.Net.TLS.Enable = true
		certBytes, err := ioutil.ReadFile("./conf/ca-cert")
		if err != nil {
			panic("kafka manager failed fo find ca-cert.file, please check directory")
		}
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			panic("kafka manager config failed to parse root certificate")
		}
		MqClusterConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            clientCertPool,
			InsecureSkipVerify: true,
		}
	}
	// 生产环境已确保Linux网络环境调用安全，所以不需要SASL校验
	if beego.AppConfig.String("runmode") == "prod" {
		MqClusterConfig.Net.SASL.Enable = false
		MqClusterConfig.Net.TLS.Enable = false
	}
	return MqClusterConfig
}

/**
  * 1.用于获取topic
 */
func (this *KafkaConfig) GetConfig() (config *sarama.Config) {
	config = saramaConfigCommon()
	return MqConfig
}

/**
  * 1.用于创建topic
 */
func (this *KafkaConfig) GetAdminClusterConfig() (config *sarama.Config) {
	config = saramaConfigCommon()
	// 这个不填会报 api version not support
	config.Version = sarama.V1_0_0_0
	config.Net.DialTimeout = time.Second * 10
	// 创建topic比较耗时，加一个最长写时间
	config.Net.WriteTimeout = time.Second * 10
	return config
}

/**
  1. 获取consumeGroups
 */
func (this *KafkaConfig) GetConsumeGroupConfig() (config *sarama.Config)  {
	config = saramaConfigCommon()
	config.Version = sarama.V0_9_0_0
	return config
}

// 公共函数，将字符串server转成切片
func (this *KafkaConfig) GetServers() []string {
	return strings.Split(beego.AppConfig.String("SERVERS"), ",")
}

// 获取消费者集群配置
func (this *KafkaConfig) GetConsumerConfig() (config *cluster.Config) {
	config = clusterConfigCommon()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

func (this *KafkaConfig) GetProducerConfig() (config *sarama.Config) {
	config = saramaConfigCommon()
	config.Producer.Return.Successes = true
	return config
}


