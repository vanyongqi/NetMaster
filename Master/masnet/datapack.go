package masnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"netMaster/Master/masiface"
	"netMaster/Master/utils"
)

//封包 拆包的模块

type DataPack struct {
	//Len  uint32
	//Type uint32
	//Data []byte
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (da *DataPack) GetHeadLen() uint32 {
	return 8
}
func (da *DataPack) Pack(msg masiface.IMessage) ([]byte, error) {
	//创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将dataLen msgID写入databuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	//将数据放入databuf中
	return dataBuff.Bytes(), nil
}

// 将包的head读出，之后根据head信息里的data长度在进行一次读
func (da *DataPack) UnPack(binaryData []byte) (masiface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	//dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen超出允许的最大包的长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large of the packet")
	}
	return msg, nil
}
