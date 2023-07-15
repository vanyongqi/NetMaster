package masiface

type Iserver interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
	//获取当前server的链接
	GetConnMgr() IConnManager

	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(conn IConnection)
	CallOnStop(connn IConnection)
}
