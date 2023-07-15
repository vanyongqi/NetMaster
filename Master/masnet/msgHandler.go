package masnet

import (
	"fmt"
	"netMaster/Master/masiface"
	"netMaster/Master/utils"
	"strconv"
)

// 消息处理模块的实现
type MsgHandle struct {
	//路由表的hash表
	Apis map[uint32]masiface.IRouter
	//业务工作worker池的数量
	WorkertPoolSize uint32
	//负责worker取任务的消息队列
	TaskQueue []chan masiface.IRequest
}

//初始化 创建MsgHandler方法

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:            make(map[uint32]masiface.IRouter),
		WorkertPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:       make([]chan masiface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启worker ，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkertPoolSize); i++ {
		//一个worker被启动·
		//1 给当前的worker对应的channel消息队列发送开辟空间第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan masiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前worker 阻塞消息从channel传递过来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan masiface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is started")
	for {
		select {
		//如果由消息过来，出列的就是一个客户端的request，并执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)

		}
	}
}

// 将消息交给TaskQueue 由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request masiface.IRequest) {
	// 将消息平均分配给不同的worker 	//分布式 可以考虑按照ip来匹配
	//根据客户端建立的ConnID来进行分配
	//轮询平均分配
	workerID := request.GetConnection().GetConnID() % mh.WorkertPoolSize
	fmt.Println("ADD ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID =", request.GetMsgId(), "to WorkerID =", workerID)
	//将数据发送给对应dworkerTaskQueue队列
	mh.TaskQueue[workerID] <- request
}
