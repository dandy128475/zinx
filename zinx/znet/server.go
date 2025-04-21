package znet

import (
	"fmt"
	"log"
	"net"
	"zinxsrc/zinx/utils"
	"zinxsrc/zinx/ziface"
)

// 定义一个服务器
type Server struct {

	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 服务器的路由
	MsgHandler ziface.IMsgHandler
	// 连接管理器
	ConnManager ziface.IConnManager
	//server销毁连接之前调用hook函数
	OnConnStop func(conn ziface.IConnection)
	// server创建连接之前调用的hook函数
	OnConnStart func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	//fmt.Println("DEBUG: Start called")
	fmt.Printf("[zinx] Server Name: %s, server ip: %s, server port : %d",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	fmt.Printf("zinxv0.1 start Listening, ip %s, Port %d\n ", s.IP, s.Port)
	//开启工作池
	s.MsgHandler.StartWorkerPool()
	// 获取tcp的Addr, 解析tcp的地址， 如果解析错误直接返回
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Printf("Failed to resolve address %s:%d: %v", s.IP, s.Port, err)
		return
	}
	var cid uint32 = 0

	// 监听服务器的地址
	Listenner, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Println("tcp listen failed", err)
		return
	}

	fmt.Println("start zinx success, ", s.Name, "now is listenning")
	// 阻塞等待客户端连接， 处理客户端请求
	for {
		conn, err := Listenner.AcceptTCP()
		if err != nil {
			log.Println("accept failed", err)
			continue
		}
		fmt.Println("new connection accepted, remote addr =", conn.RemoteAddr().String())
		if s.ConnManager.Len() > utils.GlobalObject.MaxConn {
			fmt.Println("too many connections, max conn is ", utils.GlobalObject.MaxConn)
			conn.Write([]byte("Error: Server has reached maximum connection limit\n"))
			conn.Close()
			continue
		}
		// 设置最大连接数
		//处理新业务方法与conn绑定，得到我们的连接模块
		fmt.Println("create connection start ,ID = ", cid)
		dealconn := NewConnection(s, conn, cid, s.MsgHandler)
		cid++
		fmt.Println("connection created, connID =", dealconn.GetConnID())

		// 启动当前连接的处理业务
		go dealconn.Start()

	}

}

func (s *Server) Stop() {
	fmt.Println("zinx server stop")
	s.ConnManager.ClearConn()
	fmt.Println("zinx server stop success")
}

func (s *Server) Serve() {
	go s.Start()
	select {}
}

func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgid, router)
	fmt.Println("add router success")
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

// 注册OnConnStart方法
func (s *Server) RegisterOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
	fmt.Println("register OnConnStart")
}

// 注册OnConnStop方法
func (s *Server) RegisterOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
	fmt.Println("register OnConnStop")
}

// 调用OnConnStart方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------->call OnConnStart")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Printf("------->call OnConnStop\n")
		s.OnConnStop(conn)
	}
}
