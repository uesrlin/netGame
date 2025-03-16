package siface

/**
 * @Description
 * @Date 2025/3/16 22:11
 **/

type IMsgHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgId uint32, router IRouter)
}
