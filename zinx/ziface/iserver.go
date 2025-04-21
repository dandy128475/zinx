package ziface

// 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能: 给当前服务注册一个路由方法，供客户端的连接使用
	AddRouter(msgid uint32, router IRouter)
	// 获取连接管理器
	GetConnManager() IConnManager
	// 注册OnConnStart方法
	RegisterOnConnStart(func(conn IConnection))
	// 注册OnConnStop方法
	RegisterOnConnStop(func(conn IConnection))
	// 调用OnConnStart方法
	CallOnConnStart(conn IConnection)
	// 调用OnConnStop方法
	CallOnConnStop(conn IConnection)
}
