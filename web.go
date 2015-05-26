package ksana

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
)

type webConfig struct {
	Port   int    `json:"port"`
	Cookie string `json:"cookie"`
	Expire int64  `json:"expire"`
}

func NewRouter(path string) (Router, error) {

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}
	return &router{routes: make([]Route, 0), templates: path}, nil
}

//--------------------request---------------------------------------------------

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

//-----------------------response-----------------------------------------------

type Response struct {
	path   string
	writer http.ResponseWriter
}

func (r *Response) Json(data interface{}) {
	r.writer.Header().Set("Content-Type", "application/json")
	j, e := json.Marshal(data)
	if e == nil {
		r.writer.Write(j)
	} else {
		r.Error(e)
	}

}

func (r *Response) File(req *http.Request, file string) {
	path.Join("public", file)
	http.ServeFile(r.writer, req, file)
}

func (r *Response) Xml(data interface{}) {
	r.writer.Header().Set("Content-Type", "application/xml")

	x, e := xml.MarshalIndent(data, "", "  ")
	if e == nil {
		r.writer.Write(x)
	} else {
		r.Error(e)
	}
}

func (r *Response) Text(data []byte) {
	r.writer.Write(data)
}

func (r *Response) HtmlT(file string, data interface{}) {
	t, e := template.ParseFiles(path.Join(r.path, file))
	if e != nil {
		r.Error(e)
		return
	}
	if e = t.Execute(r.writer, data); e != nil {
		r.Error(e)
	}
}

func (r *Response) TextT(buf *bytes.Buffer, tpl string, data interface{}) error {
	t, err := template.ParseFiles(path.Join(r.path, tpl))
	if err != nil {
		return err
	}
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	return nil
}

func (r *Response) Error(e error) {
	http.Error(r.writer, e.Error(), http.StatusInternalServerError)
}
