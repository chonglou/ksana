package ksana

import (
	"net/http"
	"net/url"
	"strings"
)

type PatternMux struct {
	handlers map[string][]*Handler
}

func New() *PatternMux {
	return &PatternMux{make(map[string][]*Handler)}
}

func (p *PatternMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, ph := range p.handlers[r.Method] {
		if params, ok := ph.try(r.URL.Path); ok {
			if len(params) > 0 {
				r.URL.RawQuery = url.Values(params).Encode() + "&" + r.URL.RawQuery
			}
			ph.ServeHTTP(w, r)
			return
		}
	}

	allowed := make([]string, 0, len(p.handlers))
	for meth, handlers := range p.handlers {
		if meth == r.Method {
			continue
		}

		for _, ph := range handlers {
			if _, ok := ph.try(r.URL.Path); ok {
				allowed = append(allowed, meth)
			}
		}
	}

	if len(allowed) == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Add("Allow", strings.Join(allowed, ", "))
	http.Error(w, "Method Not Allowed", 405)
}

func (p *PatternMux) Head(pat string, h http.Handler) {
	p.add("HEAD", pat, h)
}

func (p *PatternMux) Get(pat string, h http.Handler) {
	p.add("GET", pat, h)
}

func (p *PatternMux) Post(pat string, h http.Handler) {
	p.add("POST", pat, h)
}

func (p *PatternMux) Put(pat string, h http.Handler) {
	p.add("PUT", pat, h)
}

func (p *PatternMux) Patch(pat string, h http.Handler) {
	p.add("PATCH", pat, h)
}

func (p *PatternMux) Delete(pat string, h http.Handler) {
	p.add("DELETE", pat, h)
}

func (p *PatternMux) Options(pat string, h http.Handler) {
	p.add("OPTIONS", pat, h)
}

func (p *PatternMux) add(meth, pat string, h http.Handler) {
	p.handlers[meth] = append(p.handlers[meth], &Handler{pat, h})

	n := len(pat)
	if n > 0 && pat[n-1] == '/' {
		p.add(meth, pat[:n-1], http.RedirectHandler(pat, http.StatusMovedPermanently))
	}
}

func tail(pat, path string) string {
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(pat):
			if pat[len(pat)-1] == '/' {
				return path[i:]
			}
			return ""
		case pat[j] == ':':
			var nextc byte
			_, nextc, j = match(pat, isAlnum, j+1)
			_, _, i = match(path, matchPart(nextc), i)
		case path[i] == pat[j]:
			i++
			j++
		default:
			return ""
		}
	}
	return ""
}

type Handler struct {
	pat string
	http.Handler
}

func (ph *Handler) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(ph.pat):
			if ph.pat != "/" && len(ph.pat) > 0 && ph.pat[len(ph.pat)-1] == '/' {
				return p, true
			}
			return nil, false
		case ph.pat[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(ph.pat, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			p.Add(":"+name, val)
		case path[i] == ph.pat[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(ph.pat) {
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
