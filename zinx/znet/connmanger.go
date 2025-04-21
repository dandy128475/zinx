package znet

import (
	"fmt"
	"sync"
	"zinxsrc/zinx/ziface"
)

type ConnManager struct {
	// 连接ID和连接的映射
	connections map[uint32]ziface.IConnection
	// 保护连接的读写锁
	connLock sync.RWMutex
}

// NewConnManager 创建连接管理器
func NewConnManager() *ConnManager {
	cm := &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
	fmt.Println("NewConnManager: connections initialized")
	return cm
}

// Add 添加连接

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 加锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 添加连接
	cm.connections[conn.GetConnID()] = conn

	// 打印连接数量
	fmt.Printf("conn add success: %d, conn len: %d\n", conn.GetConnID(), len(cm.connections))
}

/*
func (cm *ConnManager) Add(conn ziface.IConnection) {
	if conn == nil {
		fmt.Println("ConnManager.Add: conn is nil")
		return
	}
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	if cm.connections == nil {
		fmt.Println("ConnManager.Add: connections map is nil, initializing")
		cm.connections = make(map[uint32]ziface.IConnection)
	}
	cm.connections[conn.GetConnID()] = conn
	length := len(cm.connections) // 直接计算长度，避免调用 Len()
	fmt.Printf("conn add success: %d, conn len: %d\n", conn.GetConnID(), length)
}
*/
// Remove 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 加锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除连接
	delete(cm.connections, conn.GetConnID())

	// 打印连接数量
	fmt.Printf("remove conn success: %d, conn len: \n", conn.GetConnID())
}

// Get 获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	// 获取连接
	conn, ok := cm.connections[connID]
	if !ok {
		return nil, fmt.Errorf("connection not found")
	}
	return conn, nil
}

// Len 获取连接数量
func (cm *ConnManager) Len() int {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	return len(cm.connections)
}

// ClearConn 清除所有连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 遍历所有连接，关闭连接
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("清除所有连接成功")
}
