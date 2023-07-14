package masnet

import (
	"fmt"
	"netMaster/Master/masiface"
	"strconv"
)

// 消息处理模块的实现
type MsgHandle struct {
	//路由表的hash表
	Apis map[uint32]masiface.IRouter
}

//初始化 创建MsgHandler方法

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]masiface.IRouter),
	}
}

// 执行对应路由的消息处理方法
func (mh *MsgHandle) DoMsgHandler(request masiface.IRequest) {
	//从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), " is NOT FOUND")
	}
	//根据msgID找到对应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}

// 为router添加函数绑定方式
func (mh *MsgHandle) AddRouter(msgID uint32, router masiface.IRouter) {
	//判断 当前msg绑定的API处理方法是否存在·

	if _, ok := mh.Apis[msgID]; ok {
		//id 已经注册
		//panic("repeat API  msgID:" + string(msgID))
		panic("repeat API  msgID:" + strconv.Itoa((int(msgID))))
	}
	//添加 msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("add API MsgID= ", msgID, " succ!")
}
