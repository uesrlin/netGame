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
}
