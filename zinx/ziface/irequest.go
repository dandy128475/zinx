package ziface

// 将客户端的请求连接和数据包装到request请求
type IRequest interface {
	//客户端的连接
	GetConnection() IConnection
	//客户端的数据
	GetData() []byte
	GetMsgId() uint32
	GetDataLen() uint32
}
