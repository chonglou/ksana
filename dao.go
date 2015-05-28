package ksana

import(
  "database/sql"
  "reflect"
  )

type Dao struct{
	Db *sql.DB
	Sq *Sql
	Aes *Aes
}

func (d *Dao) Load() error {
  var db, sq, aes interface{}
  var err error

  db, err = Get(reflect.TypeOf((*sql.DB)(nil)))
  if err != nil {
    return err
  }

  sq, err = Get(reflect.TypeOf((*Sql)(nil)))
  if err != nil {
    return err
  }

  aes, err = Get(reflect.TypeOf((*Aes)(nil)))
  if err != nil {
    return err
  }

  d.Db = db.(*sql.DB)
  d.Sq = sq.(*Sql)
  d.Aes = aes.(*Aes)
  return nil
}
