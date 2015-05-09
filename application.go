package ksana

import (
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"net/url"
	"os"
)

const VERSION = "v20150510"

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

type Environment struct {
	Port     int
	Mode     string
	redis    map[string]string
	database map[string]string
}

type Context struct {
	Database db
	Redis    redis
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

type Application struct {
	logger *syslog.Writer
}

func (app *Application) Abort(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
}

func (app *Application) requestHandler(resource Resource) http.HandlerFunc {
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
			app.Abort(writer, 405)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			app.Abort(writer, 500)
		}
		writer.WriteHeader(code)
		writer.Write(content)
	}
}

func (app *Application) initLogger(tag string) {
	var level syslog.Priority
	if os.Getenv("KSANA_ENVIRONMENT") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	logger, err := syslog.New(level, tag)
	if err != nil {
		Log.Fatalf("Error on init logger: %v", err)
	}
	app.logger = logger
}

func (app *Application) AddResource(resource Resource, path string) {
	http.HandleFunc(path, app.requestHandler(resource))
}

func (app *Application) Server(tag string, port int) {
	log.Printf("=> Booting Ksana %s", VERSION)
	log.Printf("=> Application starting in %s on http://0.0.0.0:%v\n", os.get, port)
	log.Println("=> Run `gails -h server` for more startup options")
	log.Println("=> Ctrl-C to shutdown server")
	app.initLogger(tag)
	//http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (app *Application) Db() {

}
func (app *Application) Redis() {

}
