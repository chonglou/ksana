package ksana_orm

type migration struct {
	Version string `json:"version"`
	Up      string `json:"up"`
	Down    string `json:"down"`
}
