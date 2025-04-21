package znet

import "zinxsrc/zinx/ziface"

type Request struct {
	//和客户端建立好的conn
	conn ziface.IConnection

	//客户端的数据
	msg ziface.IMessage
}

// 客户端的连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 客户端的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取消息ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

// 获取消息长度
func (r *Request) GetDataLen() uint32 {
	return r.msg.GetDataLen()
}
