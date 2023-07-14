package masiface

/*
封包· 拆包模块
处理粘包问题
*/

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
