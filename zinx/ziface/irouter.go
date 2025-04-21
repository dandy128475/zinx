package ziface

// 路由抽象接口
// 路由里的数据都时IRequest
type IRouter interface {
	// 处理conn业务之前的钩子函数hook
	Prehandler(request IRequest)
	// 处理conn业务的主函数hook
	Handler(request IRequest)
	// 处理conn业务之后的函数hook
	Posthandler(request IRequest)
}
