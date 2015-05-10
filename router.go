package ksana

import (
	"net/http"
	"net/url"
	"strings"
)

type HandlerFunc func(writer http.ResponseWriter, request *http.Request, context *Context)

type Router struct {
	context  *Context
	handlers map[string][]*Handler
}

func (r *Router) Init(ctx *Context) {
	r.handlers = make(map[string][]*Handler)
	r.context = ctx
}

func (r *Router) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	for _, ph := range r.handlers[req.Method] {
		if params, ok := ph.try(req.URL.Path); ok {
			if len(params) > 0 {
				req.URL.RawQuery = url.Values(params).Encode() + "&" + req.URL.RawQuery
			}
			ph.ServeHTTP(wrt, req, r.context)
			return
		}
	}

	allowed := make([]string, 0, len(r.handlers))
	for meth, handlers := range r.handlers {
		if meth == req.Method {
			continue
		}

		for _, ph := range handlers {
			if _, ok := ph.try(req.URL.Path); ok {
				allowed = append(allowed, meth)
			}
		}
	}

	if len(allowed) == 0 {
		http.NotFound(wrt, req)
		return
	}

	wrt.Header().Add("Allow", strings.Join(allowed, ", "))
	http.Error(wrt, "Method Not Allowed", 405)
}

func (r *Router) Head(p string, h *HandlerFunc) {
	r.add("HEAD", p, h)
}

func (r *Router) Get(p string, h *HandlerFunc) {
	r.add("GET", p, h)
}

func (r *Router) Post(p string, h *HandlerFunc) {
	r.add("POST", p, h)
}

func (r *Router) Put(p string, h *HandlerFunc) {
	r.add("PUT", p, h)
}

func (r *Router) Patch(p string, h *HandlerFunc) {
	r.add("PATCH", p, h)
}

func (r *Router) Delete(p string, h *HandlerFunc) {
	r.add("DELETE", p, h)
}

func (r *Router) Options(p string, h *HandlerFunc) {
	r.add("OPTIONS", p, h)
}

func (r *Router) add(meth, p string, h *HandlerFunc) {
	r.handlers[meth] = append(r.handlers[meth], &Handler{p, h})

	n := len(p)

	rf := func(writer http.ResponseWriter, request *http.Request, context *Context){
		http.RedirectHandler(p, http.StatusMovedPermanently)
	}
	if n > 0 && p[n-1] == '/' {
		r.add(meth, p[:n-1], &rf)
	}
}

type Handler struct {
	path string
	ServeHttp *HandlerFunc
}

func (h *Handler) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(h.path):
			if h.path != "/" && len(h.path) > 0 && h.path[len(h.path)-1] == '/' {
				return p, true
			}
			return nil, false
		case h.path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(h.path, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			p.Add(":"+name, val)
		case path[i] == h.path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(h.path) {
		return nil, false
	}
	return p, true
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
