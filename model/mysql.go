package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() (err error) {
	uri := "root:xxx.@tcp(47.xx.x.xxx:3306)/zxblog"
	DB, err = sql.Open("mysql", uri)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
