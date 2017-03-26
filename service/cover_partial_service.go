package service

import (
	"github.com/atongen/gosaic/model"
)

type CoverPartialService interface {
	Service
	Get(int64) (*model.CoverPartial, error)
	Insert(*model.CoverPartial) error
	BulkInsert([]*model.CoverPartial) (int64, error)
	Count(*model.Cover) (int64, error)
	Update(*model.CoverPartial) error
	Delete(*model.CoverPartial) error
	FindAll(int64, string) ([]*model.CoverPartial, error)
}
