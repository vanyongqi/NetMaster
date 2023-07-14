package masiface

type Iserver interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)
}
