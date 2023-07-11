package masnet

import (
	"fmt"
	"net"
	"netMaster/Master/masiface"
)

//连接模块

type Connection struct {
	//当前连接的socket
	Conn *net.TCPConn
	//连接id
	ConnID uint32
	//当前连接状态
	isClosed bool
	//当前连接所绑定的处理业务的方法
	handleAPI masiface.HandleFunc
	//告知当前连接已经退出/停止channel 通过管道告知要退出
	ExitChan chan bool
}

//初始化连接模块方法

func NewConnection(conn *net.TCPConn, connID uint32, callback_api masiface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// 读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("CoonID =", c.ConnID, "Reader is exit,remote add is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buffer中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//调用当前来凝结所绑定的Handle API
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", "handle is error", err)
			break
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start() CooID=", c.ConnID)
	//启动当前连接的读数据业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop().. ConnID =", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnetcion() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
