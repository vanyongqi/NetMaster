package masnet

import "netMaster/Master/masiface"

type Request struct {
	//已经和客户端建立好的链接
	conn masiface.IConnection
	msg  masiface.IMessage
	//客户端请求的数据
}

func (r *Request) GetConnection() masiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetMsgId() uint32 {
	//TODO implement me
	return r.msg.GetMsgId()
}

//func (r *Request) GetMsgLen() uint32 {
//	return r.msg.GetMsgLen()
//}
