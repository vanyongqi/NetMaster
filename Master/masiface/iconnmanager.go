package masiface

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据id获取链接
	Get(ConnID uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清除并终止所有id链接
	ClearConn()
}