package ksana

import (
	//"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
)

type HandlerFunc func(wrt http.ResponseWriter, req *http.Request, ctx *Context)

type RegexpMux struct {
	ctx      *Context
	handlers map[*regexp.Regexp]HandlerFunc
}

func (rm *RegexpMux) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	for r, h := range rm.handlers {
		if r.MatchString(url) {
			ctx.Logger.Debug("%s %s MATCH %s", req.Method, url)
			//todo 处理request 增加method
			h(wrt, req, rm.ctx)
			return
		}
	}
	http.NotFound(wrt, req)
}

func (rm *RegexpMux) add(res string, hf HandlerFunc) {
	rm.handlers[regexp.MustCompile(res)] = hf
}

func (rm *RegexpMux) routes(rs map[string]string) {
	for k, v := range rm.handlers {
		rs[k.String()] = runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	}
}
