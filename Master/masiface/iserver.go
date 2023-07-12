package masiface

type Iserver interface {
	Start()
	Stop()
	Serve()
	AddRouter(router IRouter)
}
