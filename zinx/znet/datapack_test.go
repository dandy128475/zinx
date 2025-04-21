package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 创建服务端
	Listenner, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		t.Fatalf("Failed to listen on port 8999: %v", err)
	}
	defer Listenner.Close()
	// 从客户端读取数据
	go func() {
		for {
			conn, err := Listenner.Accept()
			if err != nil {
				fmt.Println("server accept failed", err)
			}
			// 读取数据
			defer conn.Close()
			for {
				dp := NewDataPack()
				//第一次读取头
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head data error", err)
					return
				}
				// 拆包
				msgHead, err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("server unpack error", err)
					return
				}
				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*Message)
					msg.MessageData = make([]byte, msg.GetDataLen())
					// 第二次读取数据
					_, err := io.ReadFull(conn, msg.MessageData)
					if err != nil {
						fmt.Println("server unpack data error", err)
						return
					}
					fmt.Printf("recv msg id = %d, data = %s\n", msg.GetMsgId(), msg.GetData())
				}

			}

		}
	}()

	// 创建客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial failed", err)
		return
	}
	defer conn.Close()
	// 发送数据
	dp := NewDataPack()
	msg1 := &Message{
		MessageID:   1,
		MessageLen:  5,
		MessageData: []byte("hello"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack data error", err)
		return
	}
	msg2 := &Message{
		MessageID:   2,
		MessageLen:  5,
		MessageData: []byte("world"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack data error", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("client write data error", err)
		return
	}
	fmt.Println("client send data success")
	select {}
}
