package masnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"netMaster/Master/masiface"
	"netMaster/Master/utils"
	"sync"
)

// 连接模块
type Connection struct {
	//当前Conncetion隶属于的server
	TcpServer masiface.Iserver
	//当前连接的socket
	Conn *net.TCPConn
	//连接id
	ConnID uint32
	//当前连接状态
	isClosed bool
	//当前连接所绑定的处理业务的方法 ==》替换为router
	//handleAPI masiface.HandleFunc
	//告知当前连接已经退出/停止channel 通过管道告知要退出
	ExitChan chan bool //是否断开链接
	//当前链接处理的方法
	//Router masiface.IRouter
	//消息的管理MsgID和对应的处理业务的API关系
	// 用于无缓冲的管道，用于读写goroutine 之间的管道通信
	msgChan    chan []byte //消息
	MsgHandler masiface.IMsgHandle
	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块方法
func NewConnection(server masiface.Iserver, conn *net.TCPConn, connID uint32, msgHandler masiface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		//handleAPI: callback_api,//==》被替换
		isClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte), //无缓冲chan
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (s *Server) GetConnMgr() masiface.IConnManager {
	return s.ConnMgr
}

// 读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader] Goroutine is running ...")
	defer fmt.Println("CoonID =", c.ConnID, "[Reader is exit! ],remote add is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
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
		//根据绑定好的MsgID 找到对应路由的方法

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，则发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer] Goroutine is Running ")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit !]")
	//不断的阻塞等待channel消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}

		case <-c.ExitChan:
			//代表reader已经退出，此时writer也要推出
			return
		}

	}
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包 再发送

func (c *Connection) Start() {
	fmt.Println("Conn start() CooID=", c.ConnID)
	//启动当前连接的读数据业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的 创建链接之后 需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Connection Stop().. ConnID =", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.TcpServer.CallOnStop(c)
	//关闭socket连接
	c.Conn.Close()
	//回收资源
	//告知Writer关闭
	c.ExitChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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
	////将数据发送给客户端
	//if _, err := c.Conn.Write(BinaryMsg); err != nil {
	//	fmt.Println("Write msg id ", msgId, "error:", err)
	//	return errors.New("conn Write error")
	//}
	c.msgChan <- BinaryMsg
	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	//添加一个链接属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (value interface{}, err error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}

}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}
