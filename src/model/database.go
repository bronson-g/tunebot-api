package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

const driver = "mysql"
const user = "root"
const password = "Se4Q2Lp-3587"
const protocol = "tcp"
const host = "localhost"
const database = "tunebot"

func Connect() error {
	var err error
	db, err = sql.Open(driver, user+":"+password+"@"+protocol+"("+host+")/"+database)
	return err
}

func Disconnect() {
	if db != nil {
		db.Close()
		db = nil
	}
}
