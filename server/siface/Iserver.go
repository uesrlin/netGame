package siface

/**
 * @Description
 * @Date 2025/3/6 22:40
 **/

// IServer 服务器接口
type IServer interface {
	Start()                                 // 初始化资源（非阻塞）
	Stop()                                  // 回收资源
	Serve()                                 // 运行服务（阻塞）
	AddRouter(msgId uint32, router IRouter) // 添加路由
	GetManger() IconnManager                // 获取连接管理模块
	// Hook函数的注册  以函数作为形参
	SetConnStart(func(connection IConnection))
	SetConnStop(func(connection IConnection))
	//Hook函数的调用
	CallConnStart(connection IConnection)
	CallConnStop(connection IConnection)
}
