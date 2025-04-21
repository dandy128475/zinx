package znet

import (
	"fmt"
	"io"
	"net"
	"sync"
	"zinxsrc/zinx/utils"
	"zinxsrc/zinx/ziface"
)

type Connection struct {
	//当前连接的套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 连接状态(是否关闭)
	IsClosed bool
	// 连接的读写消息通道
	msgChan chan []byte
	// 判断是否退出的消息管道(由reader告知writer)
	exitChan chan bool
	// 连接的路由
	MsgHandler ziface.IMsgHandler
	// conn属于哪一个server
	TcpServer ziface.IServer
	// 连接的属性集合
	property map[string]any
	// 连接的读写锁
	propertyLock sync.RWMutex
}

func NewConnection(tcpServer ziface.IServer, conn *net.TCPConn, ConnId uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     ConnId,
		IsClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte, 10),
		exitChan:   make(chan bool, 1),
		TcpServer:  tcpServer,
		property:   make(map[string]any),
	}
	// 将当前连接添加到连接管理器中/
	c.TcpServer.GetConnManager().Add(c)

	fmt.Println("creat new connection success, connID = ", c.ConnID)
	return c
}

// 启动连接
func (c *Connection) Start() {
	fmt.Printf("conn start... connID = %d", c.ConnID)

	go c.StartReader()

	go c.StartWriter()

	//启动hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止连接
func (c *Connection) Stop() {
	fmt.Printf("conn stop... connID = %d", c.ConnID)

	if c.IsClosed {
		return
	}

	// 关闭连接
	c.IsClosed = true
	// 调用hook函数
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()
	// 关闭消息管道
	c.exitChan <- true
	// 关闭连接管理器
	if c.TcpServer.GetConnManager() != nil {
		c.TcpServer.GetConnManager().Remove(c)
	}
	// 关闭管道
	close(c.exitChan)
	close(c.msgChan)

}

// 获取当前绑定的socket conn
func (c *Connection) GetTcpConnction() *net.TCPConn {
	return c.Conn
}

// 获取连接的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的Tcp状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给远程的客户端
func (c *Connection) SendMsg(msgID uint32, msgData []byte) error {
	// 判断连接是否关闭
	if c.IsClosed {
		fmt.Println("conn is closed")
		return fmt.Errorf("conn is closed")
	}
	// 封装数据
	dp := NewDataPack()

	binarymsg, err := dp.Pack(NewMessage(msgID, msgData))
	if err != nil {
		fmt.Println("pack error", err)
		return fmt.Errorf("pack error")
	}
	// 发送数据
	c.msgChan <- binarymsg
	return nil
}

// 处理连接的读数据
func (c *Connection) StartReader() {
	fmt.Println("conn reader is started")
	defer fmt.Println("connID = ", c.ConnID, "reader is stopped, remote addr is ", c.RemoteAddr().String())

	defer c.Stop()

	for {

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		// 读取头部数据
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Println("read head data error", err)
			break
		}
		// 拆包
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack error", err)
			continue
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			// 第二次读取数据
			if _, err := io.ReadFull(c.GetTcpConnction(), data); err != nil {
				fmt.Println("server unpack data error", err)
				continue
			}
		}

		// if n == 0 {
		// 	fmt.Println("recv buf is 0")
		// 	continue
		// }
		// 处理业务
		msg.SetData(data)
		req := &Request{
			conn: c,
			msg:  msg,
		}
		fmt.Println("received msg, connID =", c.ConnID, "msgID =", msg.GetMsgId(), "data =", string(data))
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 将请求交给工作池处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 直接处理
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("conn writer is started")
	defer fmt.Println("conn ID is", c.ConnID, "Reader is exited, Address is", c.RemoteAddr().String())

	//持续读取管道中的数据，然后写数据
	for {
		select {
		// 读取消息管道中的数据
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("write data error:", err) // 记录错误但不退出
				continue                              // 继续处理下一条消息
			}
		// 读取退出管道的数据
		case <-c.exitChan:
			fmt.Println("exitChan is closed")
			return
		}
	}

}

// 设置属性
func (conn *Connection) SetProperty(key string, value any) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()
	if conn.property == nil {
		conn.property = make(map[string]any)
	}
	conn.property[key] = value
	fmt.Println("set propety success")
}

// 获取属性
func (conn *Connection) GetProperty(key string) (value any, err error) {
	conn.propertyLock.RLock()
	defer conn.propertyLock.RUnlock()
	if conn.property == nil {
		return nil, fmt.Errorf("property is nil")
	}
	if value, ok := conn.property[key]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("property not found")
}

// 删除属性
func (conn *Connection) RemoveProperty(key string) error {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()
	if conn.property == nil {
		return fmt.Errorf("property is nil")
	}
	delete(conn.property, key)
	fmt.Println("remove property success")
	return nil
}
