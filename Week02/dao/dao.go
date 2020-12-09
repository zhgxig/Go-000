package dao

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	xerror "github.com/pkg/errors"
)

var (
	userName  string = "xxxx"
	password  string = "xxxx"
	ipAddrees string = "xxxx"
	port      int    = 3306
	dbName    string = "xxxx"
	charset   string = "xxx"
)

type Info struct {
}

var SqlErrorNoRows = errors.New("sql error no rows")

var ConnectError = errors.New("connect error")


func QueryData() ([]Info, error) {
	config := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sql.Open("mysql", config)
	if err != nil {
		return nil, xerror.Wrapf(ConnectError, fmt.Sprintf("connect mysql error: %v", err))
	}

	rows, err := Db.Query("SELECT * FROM user")

	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())

		return nil, xerror.Wrapf(SqlErrorNoRows, fmt.Sprintf("sql not found error: %v", err))
	}

	for rows.Next() {
	}
	//defer rows.Close()
	var data []Info
	return data, nil
}
