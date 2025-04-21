package utils

import (
	"encoding/json"
	"os"
	"zinxsrc/zinx/ziface"
)

type GlobalObj struct {
	// Zinx的Server的名称
	Name string
	// Zinx的Server的对象
	TcpServer ziface.IServer
	// Zinx的Server的IP
	Host string
	// Zinx的Server的端口
	TcpPort int
	// Zinx的Server的版本
	Version string
	// Zinx的Server的最大连接数
	MaxConn int
	// Zinx的Server的最大数据包大小
	MaxPackageSize uint32
	// Zinx的Server的工作池的上限
	MaxWorkerTaskLen uint32
	// Zinx的Server的工作池的任务队列的长度
	WorkerPoolSize uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 解析json数据
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		TcpPort:          8999,
		Version:          "v0.4",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		Host:             "127.0.0.1",
		MaxWorkerTaskLen: 1024,
		WorkerPoolSize:   10,
	}
	// GlobalObject.Reload()
}
