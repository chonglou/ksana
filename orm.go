package ksana

type databaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ssl      string `json:"ssl"`
}

var SQL = Sql{}

type Orm struct {
}
