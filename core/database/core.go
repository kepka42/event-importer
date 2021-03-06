package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Init(conn string) error {
	db, err := sql.Open("mysql", conn+"?parseTime=true")

	if err != nil {
		return err
	}

	d.db = db
	return nil
}

func (d *Database) Close() {
	d.db.Close()
}
