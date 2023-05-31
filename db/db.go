package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Mysql *sqlx.DB

func Init(dsn string) error {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}

	Mysql = db
	return db.Ping()
}
