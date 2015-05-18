package ksana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type redisC struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Db   int    `json:"db"`
	Pool int    `json:"pool"`
}

func (rc *redisC) Url() string {
	return fmt.Sprintf("%s:%d", rc.Host, rc.Port)
}

func (rc *redisC) Shell() (string, []string) {
	return "telnet", []string{rc.Host, strconv.Itoa(rc.Port)}
}

type databaseC struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Ssl      string `json:"ssl"`
}

func (dc *databaseC) Url() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		dc.Driver, dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.Ssl)
}

func (dc *databaseC) Shell() (string, []string) {
	return "psql", []string{
		"-h", dc.Host,
		"-p", strconv.Itoa(dc.Port),
		"-d", dc.Name,
		"-U", dc.User}
}

type sessionC struct {
	Name   string `json:"name"`
	Secret []byte `json:"secret"`
}

type configuration struct {
	Name     string    `json:"name"`
	Port     int       `json:"port"`
	Env      string    `json:"env"`
	Key      []byte    `json:"key"`
	Password []byte    `json:"password"`
	Session  sessionC  `json:"session"`
	Redis    redisC    `json:"redis"`
	Database databaseC `json:"database"`
}

func writeConfig(cfg *configuration, file string) error {
	cj, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, cj, 0600)
}

func readConfig(cfg *configuration, file string) error {
	f, e := os.Open(file)
	if e != nil {
		return e
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(cfg)
}
