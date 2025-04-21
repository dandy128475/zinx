package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxsrc/zinx/utils"
	"zinxsrc/zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头长度
// | 4字节  |4字节|     |
// |消息长度|消息ID|消息体|
func (dp *DataPack) GetHeadLen() uint32 {
	return 8 // 4字节长度 + 4字节ID
}

// 封包方法：将消息数据封装成一个Message
func (dp *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	// 写入消息长度
	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetDataLen()); err != nil {
		return nil, err
	}
	// 写入消息ID
	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	// 写入消息数据
	if err := binary.Write(dataBuffer, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

// 拆包方法：将数据拆分到Message中
// 读两次
func (dp *DataPack) UnPack(binaryData []byte) (message ziface.IMessage, err error) {
	reader := bytes.NewReader(binaryData)
	// 读取消息长度
	msg := &Message{}
	if err := binary.Read(reader, binary.LittleEndian, &msg.MessageLen); err != nil {
		return nil, err
	}
	// 读取消息ID
	if err := binary.Read(reader, binary.LittleEndian, &msg.MessageID); err != nil {
		return nil, err
	}
	if (utils.GlobalObject.MaxPackageSize > 0) && (msg.MessageLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("message too large")
	}

	return msg, nil
}
