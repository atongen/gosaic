package service

import (
	"github.com/atongen/gosaic/model"
)

type GidxPartialService interface {
	Service
	Insert(*model.GidxPartial) error
	BulkInsert([]*model.GidxPartial) (int64, error)
	Update(*model.GidxPartial) error
	Delete(*model.GidxPartial) error
	Get(int64) (*model.GidxPartial, error)
	GetOneBy(string, interface{}) (*model.GidxPartial, error)
	ExistsBy(string, ...interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, ...interface{}) (int64, error)
	CountForMacro(*model.Macro) (int64, error)
	Find(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	Create(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindOrCreate(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindMissing(*model.Aspect, string, int, int) ([]*model.Gidx, error)
	CountMissing([]*model.Aspect) (int64, error)
}
