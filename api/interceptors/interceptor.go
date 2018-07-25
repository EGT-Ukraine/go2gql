package interceptors

import (
	"github.com/saturn4er/graphql"
)

type Context struct {
	Service      string
	Method       string
	Params       graphql.ResolveParams
	Request      interface{}
	PayloadError interface{}
}

type ResolveArgsInvoker func() (result interface{}, err error)
type CallMethodInvoker func(req interface{}) (result interface{}, err error)

type ResolveArgsInterceptor func(ctx *Context, next ResolveArgsInvoker) (result interface{}, err error)
type CallInterceptor func(ctx *Context, req interface{}, next CallMethodInvoker) (result interface{}, err error)

type InterceptorHandler struct {
	ResolveArgsInterceptors []ResolveArgsInterceptor
	CallInterceptors        []CallInterceptor
}

func (d *InterceptorHandler) ResolveArgs(c *Context, resolve ResolveArgsInterceptor) (res interface{}, err error) {
	chain := make([]ResolveArgsInterceptor, len(d.ResolveArgsInterceptors)+1)
	copy(chain, d.ResolveArgsInterceptors)
	chain[len(d.ResolveArgsInterceptors)] = resolve
	i := -1
	var invoker ResolveArgsInvoker
	invoker = func() (result interface{}, err error) {
		i++
		res, err := chain[i](c, invoker)
		c.Request = res
		return res, err
	}
	return invoker()
}
func (d *InterceptorHandler) Call(c *Context, req interface{}, call CallInterceptor) (res interface{}, err error) {
	chain := make([]CallInterceptor, len(d.CallInterceptors)+1)
	copy(chain, d.CallInterceptors)
	chain[len(d.CallInterceptors)] = call
	i := -1
	var invoker CallMethodInvoker
	invoker = func(req interface{}) (result interface{}, err error) {
		i++
		return chain[i](c, req, invoker)
	}
	return invoker(req)
}

func (d *InterceptorHandler) OnResolveArgs(i ResolveArgsInterceptor) {
	d.ResolveArgsInterceptors = append(d.ResolveArgsInterceptors, i)
}

func (d *InterceptorHandler) OnCall(i CallInterceptor) {
	d.CallInterceptors = append(d.CallInterceptors, i)
}
