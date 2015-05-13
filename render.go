package ksana

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path"
)

func RenderJson(wrt http.ResponseWriter, data interface{}) {
	wrt.Header().Set("Content-Type", "application/json")
	j, e := json.Marshal(data)
	if e == nil {
		wrt.Write(j)
	} else {
		http.Error(wrt, e.Error(), http.StatusInternalServerError)
	}

}

func RenderFile(wrt http.ResponseWriter, req *http.Request, file string) {
	path.Join("public", file)
	http.ServeFile(wrt, req, file)
}

func RenderXml(wrt http.ResponseWriter, data interface{}) {
	wrt.Header().Set("Content-Type", "application/xml")

	x, e := xml.MarshalIndent(data, "", "  ")
	if e == nil {
		wrt.Write(x)
	} else {
		http.Error(wrt, e.Error(), http.StatusInternalServerError)
	}

}

func RenderText(wrt http.ResponseWriter, data []byte) {
	wrt.Write(data)
}

func RenderTpl(wrt http.ResponseWriter, tpl string, data interface{}) {
	t, e := template.ParseFiles(path.Join("app/views", tpl))
	if e != nil {
		http.Error(wrt, e.Error(), http.StatusInternalServerError)
		return
	}
	if e = t.Execute(wrt, data); e != nil {
		http.Error(wrt, e.Error(), http.StatusInternalServerError)
	}

}

func Tpl2Text(buf *bytes.Buffer, tpl string, data interface{}) error {
	t, err := template.ParseFiles(path.Join("app/views", tpl))
	if err != nil {
		return err
	}
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	return nil
}
