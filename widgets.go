package ksana

type Form struct {
	Id     string  `json:"id"`
	Method string  `json:"method"`
	Action string  `json:"action"`
	Token  string  `json:"token"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Id    string      `json:"id"`
	Flag  string      `json:"type"`
	Value interface{} `json:"value"`
	Label string      `json:"label"`
}
