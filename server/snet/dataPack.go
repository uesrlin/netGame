package snet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net_game/server/siface"
)

/**
 * @Description
 * @Date 2025/3/16 10:56
 **/

type DataPack struct {
	dataLen uint32
	// 第一次构建这两个参数是没有意义的
	// 当调用Unpack方法时，会根据这两个参数来创建一个Message对象
	// 然后再调用SetData方法来设置Message对象的Data字段
	msgId uint32
	data  []byte
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + MsgId uint32(4字节)
	return 8
}

func (d *DataPack) Pack(msg siface.IMessage) ([]byte, error) {
	dataBuffer := bytes.NewBuffer([]byte{})
	// 将DataLen写入dataBuffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId写入dataBuffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将Data写入dataBuffer
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

func (d *DataPack) Unpack(data []byte) (siface.IMessage, error) {
	dataBuffer := bytes.NewReader(data)

	msg := &Message{}

	// 读取DataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读取MsgId
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if msg.DataLen > MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}
	//msg.Data = make([]byte, msg.DataLen)
	//if err := binary.Read(dataBuffer, binary.LittleEndian, msg.Data); err != nil {
	//	return nil, err
	//}

	return msg, nil
}
