package utils

import (
	"encoding/json"
	"io/ioutil"
	"netMaster/Master/masiface"
)

/*
存储一切有关框架的去哪句参数，供其他模块使用
一些参数可以通过.json文件由用户配置
*/

type GlobalObj struct {
	TcpServer masiface.Iserver
	Host      string
	TcpPort   int
	Name      string

	//
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
		Name:           "NetMaster ServerApp",
		Version:        "V0.4",
		TcpPort:        8888,
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	//从conf加载用户自定义的参数
	//GlobalObject.Reload()
}
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("Demo/netMasterV.04/conf/master.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	json.Unmarshal(data, &GlobalObject)
}
