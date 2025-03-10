package siface

/**
 * @Description
 * @Date 2025/3/9 11:49
 **/

type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
