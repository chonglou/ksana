package ksana

import (
	"http"
)

type Handler func(writer http.ResponseWriter, request *http.Request, context Context) interface{}
