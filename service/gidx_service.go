package service

import (
	"github.com/atongen/gosaic/model"
)

type GidxService interface {
	Service
	Insert(*model.Gidx) error
	Update(*model.Gidx) (int64, error)
	Delete(*model.Gidx) (int64, error)
	Get(int64) (*model.Gidx, error)
	GetOneBy(string, interface{}) (*model.Gidx, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	FindAll(string, int, int) ([]*model.Gidx, error)
}
