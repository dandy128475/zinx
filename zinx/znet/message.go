package znet

type Message struct {
	MessageID   uint32 // 消息ID
	MessageLen  uint32 // 消息长度
	MessageData []byte // 消息数据
}

// 创建一个消息包
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		MessageID:   id,
		MessageLen:  uint32(len(data)),
		MessageData: data,
	}
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.MessageID
}

// 获取消息数据
func (m *Message) GetData() []byte {
	return m.MessageData
}

// 获取消息数据长度
func (m *Message) GetDataLen() uint32 {
	return m.MessageLen
}

// 设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.MessageID = id
}

// 设置消息数据
func (m *Message) SetData(data []byte) {
	m.MessageData = data
}

// 设置消息数据长度
func (m *Message) SetDataLen(len uint32) {
	m.MessageLen = len
}
