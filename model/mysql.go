package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() (err error) {
	uri := "root:Zx200919.@tcp(47.109.17.168:3306)/zxblog"
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
