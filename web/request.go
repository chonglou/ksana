package ksana_web

import (
	"net/http"
	"net/url"
)

type Request struct {
	request *http.Request
	params  url.Values
}

func (r *Request) Path() string {
	return r.request.URL.Path
}

func (r *Request) Method() string {
	return r.request.Method
}

func (r *Request) Scheme() string {
	return r.request.URL.Scheme
}

func (r *Request) Form() url.Values {
	r.request.ParseForm()
	return r.request.Form
}

func (r *Request) Parse(rt Route) {
	names := rt.Regex().SubexpNames()
	values := rt.Regex().FindStringSubmatch(r.Path())
	for i, n := range names {
		if i > 0 {
			r.params.Set(n, values[i])
		}
	}
}

func (r *Request) Match(rt Route) bool {
	return r.Method() == rt.Method() && rt.Regex().MatchString(r.Path())
}
