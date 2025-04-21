package ziface

type IMessage interface {
	// 获取消息ID
	GetMsgId() uint32
	// 获取消息数据
	GetData() []byte
	// 获取消息数据长度
	GetDataLen() uint32
	// 设置消息ID
	SetMsgId(uint32)
	// 设置消息数据
	SetData([]byte)
	// 设置消息数据长度
	SetDataLen(uint32)
	// Pack 将消息数据封装成一个Message

}
