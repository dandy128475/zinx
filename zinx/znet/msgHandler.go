package znet

import (
	"fmt"
	"zinxsrc/zinx/utils"
	"zinxsrc/zinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32                 // 业务工作worker池的数量
	TaskQueue      []chan ziface.IRequest // 消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 执行消息处理
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " not found")
		return
	}
	// 执行注册的路由方法
	router.Prehandler(request)
	router.Handler(request)
	router.Posthandler(request)
}

// AddRouter 添加路由
func (mh *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic(fmt.Sprintf("repeated api, msgId = %d", msgId))
	}
	mh.Apis[msgId] = router
	fmt.Println("add api msgId = ", msgId, " success")
}

// 启动一个worker工作池（一个zinx框架一个工作池）
func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//当前的worker对应的消息队列
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker，阻塞等待消息队列的消息
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker
func (mh *MsgHandler) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker id = ", workerId, " is started")
	//如果有消息过来，出列的就是request，并执行当前request所绑定的业务
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(Request ziface.IRequest) {
	// 1. 获取当前要处理请求的worker id
	workerId := Request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("worker id = ", workerId, " get request msgId = ", Request.GetMsgId())
	// 2. 将消息发送给对应的worker
	mh.TaskQueue[workerId] <- Request
}
