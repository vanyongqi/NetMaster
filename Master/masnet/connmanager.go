package masnet

import (
	"errors"
	"fmt"
	"netMaster/Master/masiface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]masiface.IConnection //管理链接的集合
	connLock    sync.RWMutex                    //读写锁
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]masiface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn masiface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnMgr successfully:conn num =", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn masiface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection delete successfuly")
}

// 根据id获取链接
func (connMgr *ConnManager) Get(ConnID uint32) (masiface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[ConnID]; ok {
		//代表找到
		return conn, nil
	} else {
		return nil, errors.New("connection not Found")
	}
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有id链接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除conn并停止conn的工作

	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All connections succ! conn num = ", connMgr.Len())
}
