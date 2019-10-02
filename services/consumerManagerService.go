package services

import (
	"fmt"
	"git.mifengkong.cn/wangyu/go_micro_srv"
	"github.com/bsm/sarama-cluster"
	"github.com/pkg/errors"
	"kafka_manager/common"
	"kafka_manager/common/kafka_aliware"
	"kafka_manager/models"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	signalsMap  map[string]chan int
	signalsLock sync.Mutex
)

type ConsumerManagerService struct {
}

type ManagerConsumer struct {
	Topic         string // topic名称
	ConsumerGroup string // 消费者组名称
	Callback      string // 回调
	ConsumerId    string // 消费者唯一id
	Running       bool   // 是否正在运行
	IP            string // 为了兼容以后的分布式
}

const CLOSE_SIGNAL = 1

// 重启后，需要还原消费者现场
func (this *ConsumerManagerService) InitConsumerAndRestore() {
	resp, err := this.GetAllConsumersData()
	if err != nil {
		go_micro_srv_client.ElcLog("kafkaManager", "kafka管理系统还原现场Redis出错"+err.Error(), "error")
		panic("kafka管理系统还原现场Redis出错" + err.Error())
	}
	signalsMap = make(map[string]chan int)
	for _, consumerRecord := range resp {
		err = this.StartConsumer(consumerRecord.Topic, consumerRecord.ConsumerGroup, consumerRecord.Callback, false, consumerRecord.ConsumerId, -1)
		if err != nil {
			go_micro_srv_client.ElcLog("kafkaManager", "kafka管理系统还原现场Redis出错"+err.Error(), "error")
			panic("kafka管理系统还原现场Redis出错" + err.Error())
		}
		// 添加信号量管理，用于主动关闭某个协程
		signalsLock.Lock()
		signalsMap[consumerRecord.ConsumerId] = make(chan int, 1)
		signalsLock.Unlock()
	}
}

// 开启一个消费协程
/**
@param topic topic名称
@param consumerGroup consumerGroup消费者组名称
@param callback callback回调地址
@param newConsumer 是否为新建消费者 true: 新建 false: 还原消费者现场
@param consumerUUid 消费者协程唯一id
@param consumeTemplateId 回调模板id
*/
func (this *ConsumerManagerService) StartConsumer(topic string, consumerGroup string, callback string, newConsumer bool, consumerUUid string, consumeTemplateId int) (err error) {
	kafkaConfig := kafka_aliware.KafkaConfig{}
	consumer, err := cluster.NewConsumer(kafkaConfig.GetServers(), consumerGroup, []string{topic}, kafkaConfig.GetConsumerConfig())
	if err != nil {
		return err
	}
	// 如果停运的重启，恢复成运行状态,全量操作running为true
	if newConsumer == false {
		this.UpdateConsumerInfo(consumerUUid, true)
	}
	// 如果是新增的consumer客户端
	if newConsumer == true {
		// 唯一consumer id ，后面确保是唯一不变的
		consumerUUid = strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(1000))
		consumerBackupModel := &models.ConsumerBackup{
			Topic:             topic,
			ConsumerGroup:     consumerGroup,
			CallbackUrl:       callback,
			ConsumerId:        consumerUUid,
			Running:           true,
			ConsumeTemplateId: consumeTemplateId,
		}
		_, err := consumerBackupModel.Insert()
		if err != nil {
			return err
		}
		// 信号量map添加消费者协程id，以便信号操作停止协程
		signalsLock.Lock()
		signalsMap[consumerUUid] = make(chan int, 1)
		signalsLock.Unlock()
	}

	go func() {
		// 捕获异常
		defer func() {
			err := recover()
			// goroutine异常退出，更新当前消费者状态running为false，警报 todo 自动重启机制
			if err != nil {
				go_micro_srv_client.ElcLog("kafkaManager", "协程宕机严重错误 consumerUUid: "+consumerUUid+fmt.Sprint(err), "error")
				this.UpdateConsumerInfo(consumerUUid, false)
				consumer.Close()
			}
		}()

		for {
			select {
			case msg, more := <-consumer.Messages():
				if more {
					//fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s \n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
					// 进行回调处理
					messageProcess := &CallbackService{}
					handleError := messageProcess.ProcessMessage(callback, string(msg.Key), string(msg.Value))
					if handleError != nil {
						// 链路追踪
						common.TraceLogger(string(msg.Key), "消费者消费处理失败 "+"Topic is "+topic+" ConsumerGroup is "+consumerGroup+" consumerUUid is "+consumerUUid+handleError.Error(), "error")
					} else {
						common.TraceLogger(string(msg.Key), "消费者消费处理成功 "+"Topic is "+topic+" ConsumerGroup is "+consumerGroup+" consumerUUid is "+consumerUUid, "info")
					}
					consumer.MarkOffset(msg, "")
				}
			case ntf, more := <-consumer.Notifications():
				if more {
					// 通知，消费者重分配通知，日志记录 + 时间记录
					go_micro_srv_client.ElcLog("kafkaManager", "Topic is "+topic+" ConsumerGroup is "+consumerGroup+" consumerUUid is "+consumerUUid+fmt.Sprint(ntf), "info")
				}
			case err, more := <-consumer.Errors():
				if more {
					// 通知，消费者错误失败记录
					go_micro_srv_client.ElcLog("kafkaManager", "Topic is "+topic+" ConsumerGroup is "+consumerGroup+" consumerUUid is "+consumerUUid+err.Error(), "error")
				}
			case <-signalsMap[consumerUUid]:
				go_micro_srv_client.ElcLog("kafkaManager", consumerUUid+"主动关闭", "info")
				err = consumer.Close()
				if err != nil {
					go_micro_srv_client.ElcLog("kafkaManager", consumerUUid+"主动关闭失败"+err.Error(), "info")
				} else {
					go_micro_srv_client.ElcLog("kafkaManager", consumerUUid+"主动关闭成功", "info")
				}
				return
			}
		}
	}()
	return
}

// 通过mysql获取所有的消费者 todo 能否面向接口编程，区分redis版和mysql版
func (this *ConsumerManagerService) GetAllConsumersData() (data []*ManagerConsumer, err error) {
	consumerBackupModel := models.ConsumerBackup{}
	consumerBackupList, err := consumerBackupModel.GetConsumerBackUpList()
	if err != nil {
		return data, err
	}
	for _, v := range consumerBackupList {
		tmp := new(ManagerConsumer)
		tmp.ConsumerId = v.ConsumerId
		tmp.Callback = v.CallbackUrl
		tmp.Running = v.Running
		tmp.Topic = v.Topic
		tmp.IP = v.Ip
		tmp.ConsumerGroup = v.ConsumerGroup
		data = append(data, tmp)
	}
	return
}

// 删除某个consumer消费者
func (this *ConsumerManagerService) StopOneConsume(consumerId string) (bool, error) {
	if signalsMap[consumerId] == nil {
		return false, errors.New("当前consumerId不存在，无法删除对应的consumer")
	}
	// 要检查这个consumerId是否存在，不然会阻塞
	signalsMap[consumerId] <- CLOSE_SIGNAL
	delete(signalsMap, consumerId)

	consumerBackUpModel := models.ConsumerBackup{}
	err := consumerBackUpModel.DeleteByConsumerId(consumerId)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 去更新mysql备份数据的是否正在运行状态
func (this *ConsumerManagerService) UpdateConsumerInfo(consumerUUid string, running bool) {
	consumerBackupModel := models.ConsumerBackup{}
	data, err := consumerBackupModel.GetConsumerBackupByConsumerId(consumerUUid)
	// 如果出错或者没有数据
	if err != nil {
		return
	}
	// 更新状态就是已有状态
	if data.Running == running {
		return
	}
	data.Running = running
	data.Update()
}
