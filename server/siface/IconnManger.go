package siface

/**
 * @Description
 * @Date 2025/3/19 10:20
 **/

type IconnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connId uint32) (IConnection, error)
	Len() int
	ClearConn()
}
