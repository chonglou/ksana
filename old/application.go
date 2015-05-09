package ksana

import (
	"fmt"
)

type Application struct {
	port int
	home string
	routes map[string] interface{}
}

func (app *Application) Get(){
}

func (app *Application) Run(int port) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
