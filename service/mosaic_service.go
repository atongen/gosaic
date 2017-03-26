package service

import (
	"github.com/atongen/gosaic/model"
)

type MosaicService interface {
	Service
	Get(int64) (*model.Mosaic, error)
	Insert(*model.Mosaic) error
	Update(*model.Mosaic) (int64, error)
	GetOneBy(string, ...interface{}) (*model.Mosaic, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Mosaic, error)
}
