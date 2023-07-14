package utils

import (
	"encoding/json"
	"io/ioutil"
	"netMaster/Master/masiface"
)

/*
存储一切有关框架服务端的全局参数，供其他模块使用，
*/

type GlobalObj struct {
	TcpServer masiface.Iserver
	Host      string
	TcpPort   int
	Name      string

	Version        string //框架版本
	MaxConn        int    //最大连接数
	MaxPackageSize uint32 //数据包的最大值

}

/*
定义一个全局的对外global obj
*/
var GlobalObject *GlobalObj

func init() {
	//default value
	GlobalObject = &GlobalObj{
		Host:           "127.0.0.1",
		Name:           "NetMaster ServerApp",
		Version:        "V0.6",
		TcpPort:        8888,
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

}
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("Demo/netMasterV0.5/conf/master.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	json.Unmarshal(data, &GlobalObject)
}
