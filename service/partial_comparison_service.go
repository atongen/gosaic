package service

import (
	"github.com/atongen/gosaic/model"
)

type PartialComparisonService interface {
	Service
	Insert(*model.PartialComparison) error
	BulkInsert([]*model.PartialComparison) (int64, error)
	Update(*model.PartialComparison) error
	Delete(*model.PartialComparison) error
	DeleteBy(string, ...interface{}) error
	DeleteFrom(*model.Macro) error
	Get(int64) (*model.PartialComparison, error)
	GetOneBy(string, ...interface{}) (*model.PartialComparison, error)
	ExistsBy(string, ...interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, ...interface{}) (int64, error)
	FindAll(string, int, int, string, ...interface{}) ([]*model.PartialComparison, error)
	Find(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	Create(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	FindOrCreate(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	CountMissing(macro *model.Macro) (int64, error)
	FindMissing(*model.Macro, int) ([]*model.MacroGidxView, error)
	CreateFromView(*model.MacroGidxView) (*model.PartialComparison, error)
	GetClosest(*model.MacroPartial) (int64, error)
	GetClosestMax(*model.MacroPartial, *model.Mosaic, int) (int64, error)
	GetBestAvailable(*model.Mosaic) (*model.PartialComparison, error)
	GetBestAvailableMax(*model.Mosaic, int) (*model.PartialComparison, error)
}
