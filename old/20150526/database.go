package ksana

import (
	"database/sql"
)

type dao struct {
	db *sql.DB
}

func (d *dao) Open(driver, resource string) error {
	logger.Info("Connect to database")
	db, err := sql.Open(driver, resource)
	if err != nil {
		return err
	}
	logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return err
	}
	d.db = db
	logger.Info("Database setup successfull")
	return nil
}

// func (d *dao) Transaction(f func(tx *sql.Tx) error) {
// 	tx, err := d.db.Begin()
// 	if err != nil {
// 		logger.Err("Get transaction error: " + err.Error())
// 	}
// 	defer tx.Rollback()
// 	err = f(tx)
// 	if err != nil {
// 		return err
// 	}
// 	return tx.Commit()
// }
// func (d *dao) Call(f func(db *sql.DB) error) error {
// 	return f(d.db)
// }

// func (d *dao) Select(sq string, func(r *sql.Rows) val ...interface{}) error{
//   rs, err := d.db.Query(sq, val)
//   if err != nil{
//     return err
//   }
//   defer rs.Close()
//   for rs.Next(){
//
//   }
// }
func (d *dao) DeleteById(table string, where string, values ...interface{}) {

}

//------------------------------------------------------------------------------
var Dao = dao{}
