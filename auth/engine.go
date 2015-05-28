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
	UserName string
	Email    string
	Contact  Contact
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
	path string
}

func (ae *AuthEngine) Router(path string, r ksana.Router) {
	ae.path = path
	r.Get(ae.pattern("sign_in$"), SignInFm)
	r.Post(ae.pattern("sign_in$"), SignIn)
	r.Get(ae.pattern("register$"), RegisterFm)
	r.Post(ae.pattern("register$"), Register)
	r.Get(ae.pattern("profile$"), ProfileFm)
	r.Post(ae.pattern("profile$"), Profile)
	r.Get(ae.pattern("logs$"), Logs)
	r.Get(ae.pattern("sign_out$"), SignOut)
	r.Get(ae.pattern("unlock$"), Unlock)
	r.Get(ae.pattern("reset_password$"), ResetPassword)
	r.Get(ae.pattern("$"), Index)
	r.Get(ae.pattern("(?P<id>[\\d]+$)$"), Show)

}

func (ae *AuthEngine) pattern(path string) string {
	return fmt.Sprintf("^/%s/%s", ae.path, path)
}

func (ae *AuthEngine) Migration(mig ksana.Migrator, sq *ksana.Sql) {
	var up, idx, down bytes.Buffer
	tables := make(map[string][]string, 0)
	tables["users"] = []string{
		sq.Id(false),
		sq.String("email", false, 128, false, true, ""),
		sq.String("username", false, 128, false, false, ""),
		sq.String("token", false, 64, false, false, " "),
		sq.String("provider", false, 16, false, false, "local"),
		sq.Datetime("locked", true, ""),
		sq.Datetime("confirmed", true, ""),
		sq.Updated(),
		sq.Created()}

	fmt.Fprintf(&idx, sq.CreateIndex("", "users", true, "provider", "token"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "email"))
	fmt.Fprintf(&idx, sq.CreateIndex("", "users", false, "username"))

	tables["contacts"] = []string{
		sq.Id(false),
		sq.Int32("user_id", false, 0),
		sq.String("type", false, 16, false, false, ""),
		sq.String("val", false, 512, false, false, "")}
	fmt.Fprintf(&idx, sq.CreateIndex("", "contacts", false, "user_id", "type"))

	tables["settings"] = []string{
		sq.String("id", false, 128, false, false, ""),
		sq.Bytes("val", false, 0, true, false),
		sq.Bytes("iv", true, 32, false, false)}
	fmt.Fprintf(&idx, sq.CreateIndex("", "settings", true, "id"))

	tables["logs"] = []string{
		sq.Id(false),
		sq.String("mesage", false, 255, false, false, ""),
		sq.Created()}

	tables["roles"] = []string{
		sq.Id(false),
		sq.String("name", false, 32, false, false, ""),
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
