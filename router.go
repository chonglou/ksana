package ksana

import (
	"container/list"
	"fmt"
	"log"
	"reflect"
	"regexp"
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
	Match(url string) bool
	Parse(url string, params *Params)
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
func (r *route) Match(url string) bool {
	return r.regex.MatchString(url)
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

type Router interface {
	Get(string, ...Handler)
	Post(string, ...Handler)
	Patch(string, ...Handler)
	Put(string, ...Handler)
	Delete(string, ...Handler)
	Any(string, ...Handler)
	Resources(string, Controller)
	Add(string, string, ...Handler)
}

type router struct {
	routes *list.List
}

func (r *router) Get(pat string, hs ...Handler) {
	r.Add("GET", pat, hs)
}

func (r *router) Post(pat string, hs ...Handler) {
	r.Add("POST", pat, hs)
}

func (r *router) Patch(pat string, hs ...Handler) {
	r.Add("PATCH", pat, hs)
}

func (r *router) Put(pat string, hs ...Handler) {
	r.Add("PUT", pat, hs)
}

func (r *router) Delete(pat string, hs ...Handler) {
	r.Add("DELETE", pat, hs)
}

func (r *router) Any(pat string, hs ...Handler) {
	r.Get(pat, hs)
	r.Post(pat, hs)
	r.Patch(pat, hs)
	r.Put(pat, hs)
	r.Delete(pat, hs)
}

func (r *router) Resources(name string, ctl Controller) {
	r.Add("GET", fmt.Sprintf("/%s", name), ctl.Index)
	r.Add("GET", fmt.Sprintf("/%s/new", name), ctl.New)
	r.Add("POST", fmt.Sprintf("/%s", name), ctl.Create)
	r.Add("GET", fmt.Sprintf("/%s/(?P<id>[\\d]+)", name), ctl.Show)
	r.Add("PATCH", fmt.Sprintf("/%s/(?P<id>[\\d]+)/edit", name), ctl.Edit)
	r.Add("PATCH", fmt.Sprintf("/%s/(?P<id>[\\d]+)", name), ctl.Update)
	r.Add("PUT", fmt.Sprintf("/%s/(?P<id>[\\d]+)", name), ctl.Update)
	r.Add("DELETE", fmt.Sprintf("/%s/(?P<id>[\\d]+)", name), ctl.Destroy)
}
func (r *router) Add(mtd, pat string, hs ...Handler) {

	for _, h := range hs {
		if reflect.TypeOf(h).Kind() != reflect.Func {
			log.Fatalf("ksana handler must be a callable func")
		}
	}
	r.routes.PushBack(&route{
		method:   mtd,
		regex:    regexp.MustCompile(pat),
		handlers: hs})
}
