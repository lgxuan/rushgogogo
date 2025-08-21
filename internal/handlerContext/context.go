package handlercontext

import "github.com/elazarl/goproxy"

type HandlerContext struct {
	Context *goproxy.ProxyCtx
}

func NewHandlerContext(ctx *goproxy.ProxyCtx) *HandlerContext {
	return &HandlerContext{
		Context: ctx,
	}
}

// GetProxyCtx 获取代理上下文
func (hc *HandlerContext) GetProxyCtx() *goproxy.ProxyCtx {
	return hc.Context
}
