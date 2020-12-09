# 学习笔记

## 作业
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 思路:
dao层出现的sql.ErrNoRows, 需要wrap, 并需要打印堆栈信息, 因为上层可以针对对应的error,来做对应的处理

补充:
没有查询到数据记录对业务层来说，可能并不算错误，dao层应该提供能力使得上层能把这种情况和其它出错的情况区分开来；如果直接把原始的错误Warp后抛出，上层为了判断区分是没有查到数据还是其它错误，会不可避免得与database/sql发生耦合，这是不合理的；正确的作法是应该把ErrNoRows吞掉，替换为一个自定义的错误返回，其它错误Wrap返回;另外你的代码中也有一些问题,DB.Query返回的error， 怎么可能会是你自己定义的sqlErrNoRows呢？

### dao
```go
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
		if errors.Is(err, SqlErrorNoRows) {
			return nil, xerror.Wrap(err, "sql error no rows")
		}
		return nil, xerror.Wrapf(SqlErrorNoRows, fmt.Sprintf("sql not found error: %v", err))
	}

	for rows.Next() {
	}
	//defer rows.Close()
	var data []Info
	return data, nil
}
```

### biz
```go
package biz

import (
	"example/dao"
	"errors"
)

func Biz() ([]dao.Info, error) {
	info, err := dao.QueryData()
	if errors.Is(err, dao.SqlErrorNoRows) {
		var blank []dao.Info
		return blank, nil
	}
	if err != nil {
		var blank []dao.Info
		return blank, nil
	}
	return info, nil
}
```

### service
```go
package service

import (
	"example/biz"
	"example/dao"
)

func Service() ([]dao.Info, error) {
	return biz.Biz()
}

```

