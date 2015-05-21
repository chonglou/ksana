package ksana_web

type Config struct {
	Port   int    `json:"port"`
	Cookie string `json:"cookie"`
	Expire int64  `json:"expire"`
}
