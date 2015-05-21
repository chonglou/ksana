package ksana_web

import (
	"container/list"
	utils "github.com/chonglou/ksana/utils"
)

var logger, _ = utils.OpenLogger("ksana-web")

func New(path string) Router {
	return &router{routes: list.New(), templates: path}
}
