package service

import (
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type MosaicService interface {
	Service
}

type mosaicServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMosaicService(dbMap *gorp.DbMap) MosaicService {
	return &mosaicServiceImpl{dbMap: dbMap}
}

func (s *mosaicServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *mosaicServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Mosaic{}, "mosaics").SetKeys(true, "id")
	return nil
}
