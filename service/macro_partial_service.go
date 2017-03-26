package service

import (
	"github.com/atongen/gosaic/model"
)

type MacroPartialService interface {
	Service
	Insert(*model.MacroPartial) error
	Update(*model.MacroPartial) error
	Delete(*model.MacroPartial) error
	Get(int64) (*model.MacroPartial, error)
	GetOneBy(string, interface{}) (*model.MacroPartial, error)
	ExistsBy(string, interface{}) (bool, error)
	Count(*model.Macro) (int64, error)
	FindAll(string, int, int, string, ...interface{}) ([]*model.MacroPartial, error)
	Find(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	Create(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	FindOrCreate(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	CountMissing(*model.Macro) (int64, error)
	FindMissing(*model.Macro, string, int, int) ([]*model.CoverPartial, error)
	AspectIds(int64) ([]int64, error)
}
