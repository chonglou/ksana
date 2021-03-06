package ksana

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
)

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
	Regex() *regexp.Regexp
	Method() string
	Pattern() string
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

func (r *route) Regex() *regexp.Regexp {
	return r.regex
}

func (r *route) Pattern() string {
	return r.regex.String()
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
	templates string
	routes    []Route
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

	r.add("GET", fmt.Sprintf("%s$", name), ctl.Index)
	r.add("GET", fmt.Sprintf("%s/(?P<id>[\\d]+$)", name), ctl.Show)
	r.add("GET", fmt.Sprintf("%s/new$", name), ctl.New)
	r.add("GET", fmt.Sprintf("%s/(?P<id>[\\d]+)/edit$", name), ctl.Edit)

	r.add("POST", fmt.Sprintf("%s$", name), ctl.Create)
	r.add("PATCH", fmt.Sprintf("%s/(?P<id>[\\d]+$)", name), ctl.Update)
	r.add("PUT", fmt.Sprintf("%s/(?P<id>[\\d]+$)", name), ctl.Update)
	r.add("DELETE", fmt.Sprintf("%s/(?P<id>[\\d]+$)", name), ctl.Destroy)
}

func (r *router) add(mtd, pat string, hs []Handler) {

	logger.Debug("ROUTE ADD - " + mtd + " - " + pat)
	for _, h := range hs {
		// r.ctx.Logger.Debug(fmt.Sprintf(
		// 	"%s %s %v, %v",
		// 	mtd, pat, reflect.TypeOf(h), reflect.TypeOf(h).Kind()))
		ht := reflect.TypeOf(h)
		if ht.Kind() != reflect.Func || ht.NumOut() != 1 || ht.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			log.Fatalf("ksana handler must be a callable func")
		}

	}

	r.routes = append(r.routes, &route{
		method:   mtd,
		regex:    regexp.MustCompile(pat),
		handlers: hs})

}

func (r *router) Status(buf *bytes.Buffer) {
	for _, r := range r.routes {
		r.Status(buf)
	}
}

func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	req := Request{request: request, params: url.Values{}}
	res := Response{writer: writer, path: r.templates}

	logger.Info(fmt.Sprintf("%s %s", req.Method(), req.Path()))

	for _, rt := range r.routes {
		if req.Match(rt) {
			logger.Debug("MATCH WITH " + rt.Pattern())
			req.Parse(rt)
			err := rt.Call(func(i int, h Handler) error {
				logger.Debug(fmt.Sprintf("CALL %v %v", i, h))

				ht := reflect.TypeOf(h)
				ins := make([]reflect.Value, 0)
				for i := 0; i < ht.NumIn(); i++ {
					var pm interface{}
					ft := ht.In(i)
					switch ft {
					case reflect.TypeOf((*Request)(nil)):
						pm = &req
					case reflect.TypeOf((*Response)(nil)):
						pm = &res
					default:
						val, ok := Get(ft)
						if !ok {
							return errors.New("Unknown arg type: " + ft.String())
						}
						pm = val

					}
					ins = append(ins, reflect.ValueOf(pm))
				}
				rv := reflect.ValueOf(h).Call(ins)[0].Interface()
				if rv == nil {
					return nil
				}
				return rv.(error)
				//return h(&req, &res)

			})

			if err != nil {
				logger.Err(err.Error())
				//http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.NotFound(writer, request)
}
