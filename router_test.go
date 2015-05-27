package ksana

import (
	//"fmt"
	"errors"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

func Taaa(hf Handler) string {
	return runtime.FuncForPC(reflect.ValueOf(hf).Pointer()).Name()
	//return fmt.Sprintf("%v", reflect.TypeOf(hf))
}

func Tbbb(req *Request, res *Response) error {
	return nil
}

func TestRefletc(t *testing.T) {
	log.Printf("================REFLECT===============")

	log.Printf("Request type: " + reflect.TypeOf((*Request)(nil)).String())

	log.Printf("type of error: %s", reflect.TypeOf(errors.New("aaa")))
	f := reflect.TypeOf(Tbbb)

	log.Printf("Num In: %d, NumOut: %d", f.NumIn(), f.NumOut())
	for i := 0; i < f.NumIn(); i++ {
		log.Printf("Arg %d: %s", i, f.In(i).String())
	}
}

func TestRegexpMux(t *testing.T) {
	log.Printf("==================ROUTER=============================")
	// re := regexp.MustCompile("(?P<first>[a-zA-Z]+) (?P<last>[a-zA-Z]+)")
	// log.Println(re.MatchString("Alan Turing "))
	// log.Printf("Result: %q\n", re.SubexpNames())
	// reversed := fmt.Sprintf("AAA ${%s} ${%s}", re.SubexpNames()[2], re.SubexpNames()[1])
	// log.Println(reversed)
	// log.Println(re.ReplaceAllString("Alan Turing", reversed))

	reg := regexp.MustCompile("/users/(?P<act>[\\w]+)/(?P<id>[\\d]+$)")

	urls := []string{"/users/sdfer/12312/", "/users/sdfer/12312/qweqe", "/users/aaa/123", "/users/123/aaa", "/products/123/aaa"}

	for _, u := range urls {
		log.Printf("====== %s %s ======", u, reg.String())

		log.Printf("%v %q %s", reg.MatchString(u), reg.SubexpNames(), reg.FindStringSubmatch(u))

	}

	log.Printf(Taaa(Tbbb))
}
