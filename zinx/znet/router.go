package znet

import "zinxsrc/zinx/ziface"

// 实现router时，先嵌入baserouter这个基类，然后根据这个基类重写方法就好了
type BaseRouter struct{}

// 处理conn业务之前的钩子函数hook
func (br *BaseRouter) Prehandler(ziface.IRequest) {

}

// 处理conn业务的主函数hook
func (br *BaseRouter) Handler(ziface.IRequest) {

}

// 处理conn业务之后的函数hook
func (br *BaseRouter) Posthandler(ziface.IRequest) {

}
