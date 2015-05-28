package kuth

import (
	"database/sql"
	"github.com/chonglou/ksana"
)

func SignInFm(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	fm := ksana.NewForm("user_sign_in")
	if CurrentUser(req) == nil {
		fm.Add(ksana.Field{Id: "email", Type: "email", Label: "user.email"})
		fm.Add(ksana.Field{Id: "password", Type: "password", Label: "user.password"})
		fm.Add(ksana.Field{Id: "remember", Type: "checkbox", Value: true, Label: "user.remember"})
		fm.Submit()
		fm.Reset()
	} else {
		fm.Error("user.already_sigin")
	}
	res.Json(fm)
	return nil
}

func SignIn(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {

	return nil
}

func SignOut(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func RegisterFm(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	fm := ksana.NewForm("user_register")
	if CurrentUser(req) == nil {
		fm.Add(ksana.Field{Id: "username", Type: "text", Label: "user.username"})
		fm.Add(ksana.Field{Id: "email", Type: "email", Label: "user.email"})
		fm.Add(ksana.Field{Id: "password", Type: "password", Label: "user.password"})
		fm.Add(ksana.Field{Id: "password_confirm", Type: "password", Label: "user.password_confirm"})
		fm.Submit()
		fm.Reset()
		var agree string
		if err := Get("site.agreement", &agree); err == nil && agree != "" {
			fm.Add(ksana.Field{Id: "agree", Type: "checkbox", Value: true, Label: "user.agree", Text: agree})
		}

	} else {
		fm.Error("user.already_sigin")
	}

	res.Json(fm)
	return nil
}

func Register(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func ResetPassword(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func Unlock(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func ProfileFm(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func Profile(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func Logs(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func Show(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}

func Index(req *ksana.Request, res *ksana.Response, sq *ksana.Sql, db *sql.DB) error {
	return nil
}
