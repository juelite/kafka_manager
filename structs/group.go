package structs

type CreateConsumerGroups struct {
	//1. Consumer Group 名称以"CID_alikafka"或者"CID-alikafka"开头，以便与 Topic 进行区分
	//
	//2. 名称只能包含字母，数字，短横线(-)，下划线(_)
	//
	//3. 名称长度限制在 3-64 字节之间，长于 64 字节将被自动截取
	//
	//4. 一旦创建后不能再修改 Consumer Group 名称
	ConsumerGroup string `form:"consumer_group" valid:"Required;MinSize(3);MaxSize(64);Match(/^([a-zA-Z0-9-_]+)$/)"`
}

