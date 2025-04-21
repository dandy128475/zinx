package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"zinxsrc/zinx/znet"
)

func main() {
	fmt.Println("client0 connect starting")
	time.Sleep(1 * time.Second)

	//连接到服务端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		log.Println("client start error", err)
		return
	}
	// 循环写数据
	for {
		dp := znet.NewDataPack()
		// 封包
		msg, err := dp.Pack(znet.NewMessage(0, []byte("zinx client0 test message")))
		if err != nil {
			fmt.Println("client pack data error", err)
			return
		}
		// 发送数据
		if _, err := conn.Write(msg); err != nil {
			fmt.Println("client write data error", err)
			return
		}
		fmt.Println("client send data success")
		// 读取服务端的回显数据
		// 读取数据
		// 读取数据头
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("client read head data error", err)
			return
		}
		// 拆包
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack error", err)
			continue
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.MessageData = make([]byte, msg.GetDataLen())
			// 第二次读取数据
			if _, err := io.ReadFull(conn, msg.MessageData); err != nil {
				fmt.Println("client unpack data error", err)
				continue
			}
			fmt.Printf("recv msg id = %d, data = %s\n", msg.GetMsgId(), msg.GetData())
		}
		// 休眠1秒
		time.Sleep(1 * time.Second)

	}

}
