package ziface

import (
	"net"
)

type IConnection interface {
	// 启动连接
	Start()
	// 停止连接
	Stop()
	// 获取当前绑定的socket conn
	GetTcpConnction() *net.TCPConn
	// 获取连接的连接ID
	GetConnID() uint32
	// 获取远程客户端的Tcp状态
	RemoteAddr() net.Addr
	// 发送数据给远程的客户端
	SendMsg(uint32, []byte) error
	// 设置属性
	SetProperty(key string, value any)
	// 获取属性
	GetProperty(key string) (value any, err error)
	// 删除属性
	RemoveProperty(key string) error
}

// 定义一个处理连接业务的方法
type HandlerFunc func(*net.TCPConn, []byte, int) error
