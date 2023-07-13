package masnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 支负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/

	// 1 创建socektTCP  server
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("sercer listener err", err)
		return
	}
	//创建一个go承载，负责从客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端请求，拆包的过程
				dp := NewDataPack()
				for {
					//1 conn都len id
					//2 conn读取data
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHEAD, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpacke err", err)
						return
					}
					//第二次读
					if msgHEAD.GetMsgLen() > 0 {
						//,asg是有数据的，进行第二次读取数据
						msg := msgHEAD.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中，读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}
						//消息读取完毕
						fmt.Println("---> resv Meg ID ", msg.Id, " --datalen ", msg.DataLen, "--data", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	// 2 从客户端读取数据，进行拆包处理

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("err")
		return
	}
	dp := NewDataPack()
	//模拟粘包过程，封装两个msg一起发送
	//封装第一个msg 包1
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'f', 'y', 'q', 'l', 'o'},
	}
	send1, err := dp.Pack(msg1)
	if err != nil {
		return
	}
	//封装第二个msg 包2
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'v', 'e', 'z', 'm', 't', 'y', 'a'},
	}
	send2, err := dp.Pack(msg2)
	if err != nil {
		return
	}
	//将两个包粘在一起
	send1 = append(send1, send2...)
	//一次性发送粘包
	conn.Write(send1)
	//client 阻塞
	select {}
}
