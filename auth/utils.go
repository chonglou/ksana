package kuth

import (
	"database/sql"
	"github.com/chonglou/ksana"
"errors"
)


func CurrentUser(req *ksana.Request) *User {
	return nil
}

type Dao struct{
	 ksana.Dao
}

func (d *Dao) Get(key string, val interface{}, encrypt bool) error {

	rs, err := d.Db.Query(d.Sq.Select(
		"settings",
		[]string{"val", "iv"},
		"id = $1","", 0, 1), key)
	if err != nil {
		return err
	}
	defer rs.Close()
	if rs.Next() {
		var bs, iv []byte
		err = rs.Scan(&bs, &iv)
		if err != nil {
			return err
		}
		if encrypt {
			bs = d.Aes.Decrypt(bs, iv)
		}
		return ksana.Bit2obj(bs, val)
	}
	return errors.New("Not exist!")
}

func (d *Dao) Set(key string, val interface{}, encrypt bool) error {
	bs, err := ksana.Obj2bit(val)
	if err != nil {
		return err
	}

	var iv []byte
	if encrypt {
		bs, iv = d.Aes.Encrypt(bs)
	}


	var rs *sql.Rows

	rs, err = d.Db.Query(d.Sq.Count("settings", "id = $1"), key)
	if err != nil {
		return err
	}
	defer rs.Close()
	rs.Next()
	var count int
	err = rs.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = d.Db.Exec(d.Sq.Insert("settings", []string{"id", "val", "iv"}), key, bs, iv)
	} else {
		_, err = d.Db.Exec(d.Sq.Update("settings", "val=$1, iv=$2", "id=$3"), bs, iv, key)
	}

	return err
}

var dao= Dao{}
	var logger,_ = ksana.OpenLogger("ksana-auth")

func init(){
		if err:=dao.Load(); err != nil{
			logger.Err("Failed at load bean: %s"+ err.Error())
		}
}
