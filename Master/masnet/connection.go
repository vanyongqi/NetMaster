package masnet

import (
	"errors"
	"fmt"
	"io"
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
	//当前连接所绑定的处理业务的方法 ==》替换为router
	//handleAPI masiface.HandleFunc
	//告知当前连接已经退出/停止channel 通过管道告知要退出
	ExitChan chan bool
	//当前链接处理的方法
	Router masiface.IRouter
}

// 初始化连接模块方法
func NewConnection(conn *net.TCPConn, connID uint32, router masiface.IRouter) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		//handleAPI: callback_api,//==》被替换
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("CoonID =", c.ConnID, "Reader is exit,remote add is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buffer中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建一个拆包解包的对象工具
		dp := NewDataPack()
		//读取客户端的MsgHead 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg header error", err)
			break
		}
		//拆包 得到msgID 和msgDataLen，放到消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据datalen，再次读取data，放到msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen()) //拆包是为了获得包的长度，然后进行读取
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		//req是一个对象
		req := Request{
			conn: c,
			msg:  msg,
		}
		//c.Router.PreHandle(&req)
		go func(request masiface.IRequest) {
			//	c.Router.PreHandle(request)
			c.Router.Handle(request) //c.Router.Handle(request)
			//	c.Router.PostHandle(request)
		}(&req)

		//调用路由，从路由中找到方法
		//if err:= c.HandleAPI(c.Conn,buf,cnt);err!=nil{
		//fmt.Println("ConnID",c.ConnID,"handle is error",err)
		//break
		//}

	}
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包 再发送

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

func (c *Connection) GetTCPConnection() *net.TCPConn {
	//TODO implement me
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection Closed When send error")
	}
	//将data进行封包 MsgDataLen /MsgID / Data
	dp := NewDataPack()
	BinaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println(" pack message id error ", msgId)
	}
	//将数据发送给客户端
	if _, err := c.Conn.Write(BinaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, "error:", err)
		return errors.New("conn Write error")
	}
	return nil
}
