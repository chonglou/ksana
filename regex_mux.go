package ksana

import (
	//"fmt"
	"net/http"
	"regexp"
)

type RegexMux struct {
	handlers map[string][]*Handler
}

func (p *RegexMux) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {

	http.NotFound(wrt, req)
}
