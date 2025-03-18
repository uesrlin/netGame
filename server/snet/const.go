package snet

/**
 * @Description
 * @Date 2025/3/16 21:33
 **/

const (
	MaxConn          uint32 = 1000  // 最大连接数
	MaxPackageSize   uint32 = 4096  // 数据包最大值
	WorkerPoolSize   uint32 = 10    // 线程池大小
	MaxWorkerTaskLen uint32 = 500   // 消息队列数量
	ErrWorkerId      int    = -9999 // 线程池错误
)
