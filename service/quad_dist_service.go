package service

import (
	"github.com/atongen/gosaic/model"
)

type QuadDistService interface {
	Service
	Get(int64) (*model.QuadDist, error)
	Insert(*model.QuadDist) error
	GetWorst(*model.Macro, int, int) (*model.CoverPartialQuadView, error)
}
