package ksana

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"net/http"
	"net/url"
	"os"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(values url.Values) (int, interface{}) {
	return 405, ""
}

func (PostNotSupported) Post(values url.Values) (int, interface{}) {
	return 405, ""
}

func (PutNotSupported) Put(values url.Values) (int, interface{}) {
	return 405, ""
}

func (DeleteNotSupported) Delete(values url.Values) (int, interface{}) {
	return 405, ""
}

type API struct {
	logger *syslog.Writer
}

func (api *API) Abort(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var data interface{}
		var code int

		request.ParseForm()
		method := request.Method
		values := request.Form

		switch method {
		case GET:
			code, data = resource.Get(values)
		case POST:
			code, data = resource.Post(values)
		case PUT:
			code, data = resource.Put(values)
		case DELETE:
			code, data = resource.Delete(values)
		default:
			api.Abort(writer, 405)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			api.Abort(writer, 500)
		}
		writer.WriteHeader(code)
		writer.Write(content)
	}
}

func (api *API) InitLogger(tag string) {
	var level int
	if os.Getenv("KSANA_ENVIRONMENT") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	api.logger = syslog.New(level, tag)
}

func (api *API) AddResource(resource Resource, path string) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
