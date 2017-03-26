package service

import (
	"github.com/atongen/gosaic/model"
)

type AspectService interface {
	Service
	Insert(*model.Aspect) error
	Get(int64) (*model.Aspect, error)
	Count() (int64, error)
	Find(int, int) (*model.Aspect, error)
	Create(int, int) (*model.Aspect, error)
	FindOrCreate(int, int) (*model.Aspect, error)
	FindIn([]int64) ([]*model.Aspect, error)
}
