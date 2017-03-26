package service

import (
	"github.com/atongen/gosaic/model"
)

type MacroService interface {
	Service
	Get(int64) (*model.Macro, error)
	Insert(*model.Macro) error
	Update(*model.Macro) error
	Delete(*model.Macro) error
	GetOneBy(string, ...interface{}) (*model.Macro, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Macro, error)
}
