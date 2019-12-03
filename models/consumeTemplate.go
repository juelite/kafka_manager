package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// table name : topic_callback_template
type TopicCallbackTemplate struct {
	Id          int       `json:"id"`
	Topic       string    `json:"topic"`
	CallbackUrl string    `json:"callback_url"`
	CreatedTime time.Time `json:"created_time"`
}

func init() {
	orm.RegisterModel(new(TopicCallbackTemplate))
}

// 创建一条新consumerTemplate数据
func (this *TopicCallbackTemplate) Insert() (int64, error) {
	return orm.NewOrm().Insert(this)
}

// 查询topic和callback_url模板是否已存在
func (this *TopicCallbackTemplate) IsExistTopicAndCallback() (bool, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable(this).Filter("topic", this.Topic).Filter("callback_url", this.CallbackUrl).Count()
	return count > 0, err
}

// 查询所有的topic callback 模板
func (this *TopicCallbackTemplate) GetAllTopicCallbacks() (res []*TopicCallbackTemplate, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(this).All(&res)
	return
}

// 通过id查询topic回调模板
func (this *TopicCallbackTemplate) GetTopicCallBackById() (data TopicCallbackTemplate, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(this).Filter("id", this.Id).One(&data)
	return
}

// 删除指定的模板
func (this *TopicCallbackTemplate) DeleteByTemplateId(templateId int) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(this).Filter("id", templateId).Delete()
	return
}
