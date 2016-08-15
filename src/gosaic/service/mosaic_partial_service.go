package service

import (
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type MosaicPartialService interface {
	Service
}

type mosaicPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMosaicPartialService(dbMap *gorp.DbMap) MosaicPartialService {
	return &mosaicPartialServiceImpl{dbMap: dbMap}
}

func (s *mosaicPartialServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *mosaicPartialServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.MosaicPartial{}, "mosaic_partials").SetKeys(true, "id")
	return nil
}
