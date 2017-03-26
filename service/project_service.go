package service

import (
	"github.com/atongen/gosaic/model"
)

type ProjectService interface {
	Service
	Get(int64) (*model.Project, error)
	Insert(*model.Project) error
	Update(*model.Project) (int64, error)
	GetOneBy(string, ...interface{}) (*model.Project, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Project, error)
}
