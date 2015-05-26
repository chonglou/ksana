package ksana

import (
	"container/list"
	utils "github.com/chonglou/ksana/utils"
)

type webConfig struct {
	Port   int    `json:"port"`
	Cookie string `json:"cookie"`
	Expire int64  `json:"expire"`
}

func NewWeb(path string) Router {
	return &router{routes: list.New(), templates: path}
}
