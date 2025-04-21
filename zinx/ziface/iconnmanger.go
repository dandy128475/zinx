package ziface

type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 删除连接
	Remove(conn IConnection)
	// Get 获取连接
	Get(connID uint32) (IConnection, error)
	// Len 获取连接数量
	Len() int
	// ClearConn 清除所有连接
	ClearConn()
}
