package kuth

import (
	"database/sql"
	"github.com/chonglou/ksana"
)

func CurrentUser(req *ksana.Request) *User {
	return nil
}

func Get(key string, val interface{}) error {

	return nil
}

func Set(key string, val interface{}, encrypt bool) error {
	db, err := ksana.Get(&sql.DB{})
	if err != nil {
		return err
	}
	if encrypt {

		if err != nil {
			return err
		}
	}
}
