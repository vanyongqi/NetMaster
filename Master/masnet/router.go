package masnet

import "netMaster/Master/masiface"

// 实现BaseRouter ，嵌入这个基类，然后根据需要对这
// 也就是说用户实现的时候，不需要去写多余的Pre 或者PostHandle
type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request masiface.IRequest) {

}

func (br *BaseRouter) Handle(request masiface.IRequest) {

}

func (br *BaseRouter) PostHandle(request masiface.IRequest) {

}
