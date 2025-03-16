package snet

/**
 * @Description
 * @Date 2025/3/16 10:42
 **/

type Message struct {
	DataLen uint32
	Id      uint32
	Data    []byte
}

func NewMsgPackage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(u uint32) {
	m.Id = u
}

func (m *Message) SetMsgLen(u uint32) {
	m.DataLen = u
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
