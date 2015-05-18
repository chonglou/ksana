package ksana

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path"
)

type Response struct {
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
	t, e := template.ParseFiles(path.Join("app/views", file))
	if e != nil {
		r.Error(e)
		return
	}
	if e = t.Execute(r.writer, data); e != nil {
		r.Error(e)
	}
}

func (r *Response) TextT(buf *bytes.Buffer, tpl string, data interface{}) error {
	t, err := template.ParseFiles(path.Join("app/views", tpl))
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
