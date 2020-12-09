package service

import (
	"example/biz"
	"example/dao"
)

func Service() ([]dao.Info, error) {
	return biz.Biz()
}
