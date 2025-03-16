package siface

/**
 * @Description
 * @Date 2025/3/8 12:28
 **/

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgId() uint32
}
