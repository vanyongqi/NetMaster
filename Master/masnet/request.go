package masnet

import "netMaster/Master/masiface"

type Request struct {
	//已经和客户端建立好的链接
	conn masiface.IConnection
	data []byte
	//客户端请求的数据
}

func (r *Request) GetConnection() masiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
