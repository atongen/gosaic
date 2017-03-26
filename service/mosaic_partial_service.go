package service

import (
	"github.com/atongen/gosaic/model"
)

type MosaicPartialService interface {
	Service
	Get(int64) (*model.MosaicPartial, error)
	Insert(*model.MosaicPartial) error
	Count(*model.Mosaic) (int64, error)
	CountMissing(*model.Mosaic) (int64, error)
	GetMissing(*model.Mosaic) (*model.MacroPartial, error)
	GetRandomMissing(*model.Mosaic) (*model.MacroPartial, error)
	FindAllPartialViews(*model.Mosaic, string, int, int) ([]*model.MosaicPartialView, error)
	FindRepeats(*model.Mosaic, int) ([]int64, error)
}
