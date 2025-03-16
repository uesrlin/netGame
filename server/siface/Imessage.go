package siface

/**
 * @Description
 * @Date 2025/3/16 10:42
 **/

type IMessage interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte
	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}
