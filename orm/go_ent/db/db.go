package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xiaozefeng/go-example/orm/go_ent/ent"
)

func GetClient() (*ent.Client, error) {
	var (
		username = "root"
		password = "1qaz@WSX"
		host     = "127.0.0.1"
		port     = 3306
		database = "foo"
	)
	var dbURL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", username,
		password,
		host,
		port,
		database)
	return ent.Open("mysql", dbURL)
}
