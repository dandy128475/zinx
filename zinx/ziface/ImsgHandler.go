package ziface

type IMsgHandler interface {
	// DoMsgHandler 执行消息处理
	DoMsgHandler(request IRequest)
	// AddRouter 添加路由
	AddRouter(msgId uint32, router IRouter)

	StartWorkerPool() //启动工作池

	SendMsgToTaskQueue(request IRequest) //将消息发送到任务队列
}
