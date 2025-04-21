package main

import (
	"fmt"
	"zinxsrc/zinx/ziface"
	"zinxsrc/zinx/znet"
)

type pingrouter struct {
	znet.BaseRouter
}
type hello struct {
	znet.BaseRouter
}

func (p *pingrouter) Handler(request ziface.IRequest) {
	fmt.Println("call router Handler")
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping...")); err != nil {
		fmt.Println("send msg failed, err = ", err)
	} else {
		fmt.Println("send msg success")
	}

}
func (h *hello) Handler(request ziface.IRequest) {
	fmt.Println("call router Handler")
	fmt.Println("recv from client: msgId = ", request.GetMsgId(), ", data = ", string(request.GetData()))

	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		fmt.Println("send msg failed, err = ", err)
	} else {
		fmt.Println("send msg success")
	}

}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin")
	if err := conn.SendMsg(0, []byte("DoConnectionBegin")); err != nil {
		fmt.Println("send msg failed, err = ", err)
	} else {
		fmt.Println("send msg success")
	}

	// 设置连接属性
	conn.SetProperty("name", "zinx")
	conn.SetProperty("age", 18)

}

func DoConnectionAfter(conn ziface.IConnection) {
	fmt.Println("DoConnectionAfter")
	if err := conn.SendMsg(1, []byte("DoConnectionAfter")); err != nil {
		fmt.Println("send msg failed, err = ", err)
	} else {
		fmt.Println("send msg success")
	}
	// 获取连接属性
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("name = ", name)
	} else {
		fmt.Println("get name failed, err = ", err)
	}
	if age, err := conn.GetProperty("age"); err == nil {
		fmt.Println("age = ", age)
	} else {
		fmt.Println("get age failed, err = ", err)
	}
}

// 基于zinx开发的服务端应用程序
func main() {
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer("zinx v0.6")

	//添加一个自定义router
	s.AddRouter(0, &pingrouter{})
	s.AddRouter(1, &hello{})

	//注册OnConnStart和OnConnStop钩子方法
	s.RegisterOnConnStart(DoConnectionBegin)
	s.RegisterOnConnStop(DoConnectionAfter)

	//启动server
	s.Serve()
}
