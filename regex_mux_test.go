package ksana

import (
	//"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

func TestRouter(t *testing.T) {
	var rt Route
	rt = &route{}
	var rtr Router
	rtr = &router{}

	log.Printf("ROUTE: %s %v", rt.Method(), rtr)
}

func Taaa(hf HandlerFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(hf).Pointer()).Name()
	//return fmt.Sprintf("%v", reflect.TypeOf(hf))
}

func Tbbb(wrt http.ResponseWriter, req *http.Request, ctx *Context) {

}

func TestRegexpMux(t *testing.T) {
	// re := regexp.MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
	// log.Println(re.MatchString("Alan Turing "))
	// log.Printf("Result: %q\n", re.SubexpNames())
	// reversed := fmt.Sprintf("AAA ${%s} ${%s}", re.SubexpNames()[2], re.SubexpNames()[1])
	// log.Println(reversed)
	// log.Println(re.ReplaceAllString("Alan Turing", reversed))

	reg := regexp.MustCompile("/users/(?P<act>[\\w]+)/(?P<id>[\\d]+)")

	urls := []string{"/users/aaa/123", "/users/123/aaa", "/products/123/aaa"}

	for _, u := range urls {
		log.Printf("====== %s %s ======", u, reg.String())

		log.Printf("%v %q %s", reg.MatchString(u), reg.SubexpNames(), reg.FindStringSubmatch(u))

	}

	log.Printf(Taaa(Tbbb))
}
