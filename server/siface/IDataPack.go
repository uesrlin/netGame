package siface

/**
 * @Description
 * @Date 2025/3/16 10:56
 **/

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
