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
