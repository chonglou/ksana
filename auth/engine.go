package kuth

import (
	"bytes"
	"fmt"
	"github.com/chonglou/ksana"
	"time"
)

type Contact struct {
	Qq       string
	Wechat   string
	Skype    string
	Linkedin string
	Factbook string
	Logo     string
}

type User struct {
	FirstName string
	LastName  string
	Email     string
	Contact   Contact
}

type Setting struct {
	Key string
	Val string
}

type Log struct {
	Id      int
	Message string
	Created time.Time
}

type Role struct {
	Id      int
	Name    string
	Rid     int
	Rtype   string
	Created time.Time
	Updated time.Time
}

type AuthEngine struct {
}

func (ae *AuthEngine) Router(path string, r ksana.Router) {
	r.Resources(fmt.Sprintf("%s/users", path), ksana.Controller{
		Index:   []ksana.Handler{},
		Show:    []ksana.Handler{},
		New:     []ksana.Handler{},
		Create:  []ksana.Handler{},
		Edit:    []ksana.Handler{},
		Update:  []ksana.Handler{},
		Destroy: []ksana.Handler{},
	})
}

func (ae *AuthEngine) Migration(mig ksana.Migrator, sq *ksana.Sql) {
	var up, idx, down bytes.Buffer
	tables := make(map[string][]string, 0)
	tables["users"] = []string{
		sq.Id(false),
		sq.String("email", false, 127, false, true, ""),
		sq.String("first_name", false, 31, false, false, ""),
		sq.String("middle_name", false, 31, false, true, ""),
		sq.String("last_name", false, 31, false, false, ""),
		sq.String("token", false, 63, false, false, " "),
		sq.String("provider", false, 15, false, false, "local"),
		sq.Datetime("locked", true, ""),
		sq.Datetime("confirmed", true, ""),
		sq.Updated(),
		sq.Created()}

	fmt.Fprintf(&idx, sq.CreateIndex("", "users", true, "provider", "token"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "email"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "first_name"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "last_name"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "middle_name"))

	tables["contacts"] = []string{
		sq.Id(false),
		sq.Int32("user_id", false, 0),
		sq.String("type", false, 15, false, false, ""),
		sq.String("val", false, 511, false, false, "")}
	fmt.Fprintf(&idx, sq.CreateIndex("", "contacts", false, "user_id", "type"))

	tables["settings"] = []string{
		sq.String("id", false, 127, false, false, ""),
		sq.Bytes("val", false, 0, true, false),
		sq.Bytes("iv", true, 32, false, false)}
	fmt.Fprintf(&idx, sq.CreateIndex("", "settings", true, "id"))

	tables["logs"] = []string{
		sq.Id(false),
		sq.String("mesage", false, 255, false, false, ""),
		sq.Created()}

	tables["roles"] = []string{
		sq.Id(false),
		sq.String("name", false, 31, false, false, ""),
		sq.Int32("r_id", false, 0),
		sq.String("r_type", false, 255, false, false, ""),
		sq.Date("startup", false, "now"),
		sq.Date("shutdown", false, "9999-12-31"),
		sq.Created()}

	fmt.Fprintf(&idx, sq.CreateIndex("", "roles", true, "name", "r_id", "r_type"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "roles", false, "name"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "roles", false, "r_type"))

	tables["users_roles"] = []string{
		sq.Int32("user_id", false, 0),
		sq.Int32("role_id", false, 0)}
	fmt.Fprintf(&idx, sq.CreateIndex("", "users_roles", true, "user_id", "role_id"))

	for k, v := range tables {
		fmt.Fprintf(&up, sq.CreateTable(k, v...))
		fmt.Fprintf(&down, sq.DropTable(k))
	}

	up.Write(idx.Bytes())
	mig.Add("20150526", "ksana_init", up.String(), down.String())
}
