package ksana

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
)

type Params map[string]string

type Handler interface{}

type Controller struct {
	Index   []Handler
	New     []Handler
	Create  []Handler
	Show    []Handler
	Edit    []Handler
	Update  []Handler
	Destroy []Handler
}

type Route interface {
	Method() string
	Pattern() string
	Match(method string, url string) bool
	Parse(url string, params Params)

	Status(buf *bytes.Buffer)
	Call(func(i int, h Handler) error) error
}

type route struct {
	method   string
	regex    *regexp.Regexp
	handlers []Handler
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Pattern() string {
	return r.regex.String()
}

func (r *route) Match(mtd string, url string) bool {
	return mtd == r.method && r.regex.MatchString(url)
}

func (r *route) Parse(url string, params Params) {
	names := r.regex.SubexpNames()
	values := r.regex.FindStringSubmatch(url)
	for i, n := range names {
		if i > 0 {
			params[n] = values[i]
		}
	}
}

func (r *route) Status(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "===== %s %s =====\n", r.method, r.Pattern())
	for i, h := range r.handlers {
		fmt.Fprintf(
			buf,
			"%d: %s %v\n",
			i+1,
			runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name(),
			reflect.TypeOf(h))
	}
}

func (r *route) Call(f func(int, Handler) error) error {
	for i, h := range r.handlers {
		if err := f(i, h); err != nil {
			return err
		}
	}
	return nil
}

type Router interface {
	Get(string, ...Handler)
	Post(string, ...Handler)
	Patch(string, ...Handler)
	Put(string, ...Handler)
	Delete(string, ...Handler)
	Any(string, ...Handler)
	Resources(string, Controller)

	Status(*bytes.Buffer)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type router struct {
	ctx    *Context
	routes *list.List
}

func (r *router) Get(pat string, hs ...Handler) {
	r.add("GET", pat, hs)
}

func (r *router) Post(pat string, hs ...Handler) {
	r.add("POST", pat, hs)
}

func (r *router) Patch(pat string, hs ...Handler) {
	r.add("PATCH", pat, hs)
}

func (r *router) Put(pat string, hs ...Handler) {
	r.add("PUT", pat, hs)
}

func (r *router) Delete(pat string, hs ...Handler) {
	r.add("DELETE", pat, hs)
}

func (r *router) Any(pat string, hs ...Handler) {
	r.add("GET", pat, hs)
	r.add("POST", pat, hs)
	r.add("PATCH", pat, hs)
	r.add("PUT", pat, hs)
	r.add("DELETE", pat, hs)
}

func (r *router) Resources(name string, ctl Controller) {

	r.add("GET", fmt.Sprintf("/%s$", name), ctl.Index)
	r.add("GET", fmt.Sprintf("/%s/(?P<id>[\\d]+$)", name), ctl.Show)
	r.add("GET", fmt.Sprintf("/%s/new$", name), ctl.New)
	r.add("GET", fmt.Sprintf("/%s/(?P<id>[\\d]+)/edit$", name), ctl.Edit)

	r.add("POST", fmt.Sprintf("/%s$", name), ctl.Create)
	r.add("PATCH", fmt.Sprintf("/%s/(?P<id>[\\d]+$)", name), ctl.Update)
	r.add("PUT", fmt.Sprintf("/%s/(?P<id>[\\d]+$)", name), ctl.Update)
	r.add("DELETE", fmt.Sprintf("/%s/(?P<id>[\\d]+$)", name), ctl.Destroy)
}

func (r *router) add(mtd, pat string, hs []Handler) {

	r.ctx.Logger.Debug("ROUTE ADD - " + mtd + " - " + pat)
	for _, h := range hs {
		// r.ctx.Logger.Debug(fmt.Sprintf(
		// 	"%s %s %v, %v",
		// 	mtd, pat, reflect.TypeOf(h), reflect.TypeOf(h).Kind()))
		if reflect.TypeOf(h).Kind() != reflect.Func {
			log.Fatalf("ksana handler must be a callable func")
		}
	}

	r.routes.PushBack(&route{
		method:   mtd,
		regex:    regexp.MustCompile(pat),
		handlers: hs})
}

func (r *router) Status(buf *bytes.Buffer) {
	for it := r.routes.Front(); it != nil; it = it.Next() {
		it.Value.(Route).Status(buf)
	}
}

func (r *router) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	url, method := req.URL.Path, req.Method

	r.ctx.Logger.Info(fmt.Sprintf("%s %s", method, url))

	for it := r.routes.Front(); it != nil; it = it.Next() {
		rt := it.Value.(Route)
		if rt.Match(method, url) {
			//r.ctx.Logger.Debug(fmt.Sprintf("MATCH WITH %s", rt.Pattern()))
			err := rt.Call(func(i int, h Handler) error {
				//todo 处理
				r.ctx.Logger.Debug(fmt.Sprintf("%v %v", i, h))
				return nil
			})

			if err != nil {
				http.Error(wrt, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}
	http.NotFound(wrt, req)
}
