package service

import (
	"github.com/atongen/gosaic/model"
)

type CoverService interface {
	Service
	Get(int64) (*model.Cover, error)
	Insert(*model.Cover) error
	Update(*model.Cover) error
	Delete(*model.Cover) error
	GetOneBy(string, ...interface{}) (*model.Cover, error)
	FindAll(string) ([]*model.Cover, error)
}
