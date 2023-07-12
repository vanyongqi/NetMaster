package masiface

// 指令 -> 数据处理方式 方法
// 不同的消息对应不同的处理方式
type IRouter interface {
	//设置“钩子”，应用程序可以在系统级对所有消息、
	//事件进行过滤，访问在正常情况下无法访问的消息
	//处理conn 之前的方法 前分发 hook
	PreHandle(request IRequest)
	//处理conn的业务方法
	Handle(request IRequest)
	//处理conn业务之后的方法
	PostHandle(request IRequest)
}
