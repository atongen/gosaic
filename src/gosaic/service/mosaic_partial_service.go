package service

import (
	"errors"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type MosaicPartialService interface {
	Service
	Get(int64) (*model.MosaicPartial, error)
	Insert(*model.MosaicPartial) error
	GetRandomMissing(int64) *model.MacroPartial
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

func (s *mosaicPartialServiceImpl) Get(id int64) (*model.MosaicPartial, error) {
	c, err := s.DbMap().Get(model.MosaicPartial{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	m, ok := c.(*model.MosaicPartial)
	if !ok {
		return nil, errors.New("Unable to type cast mosaic partial")
	}

	if m.Id == int64(0) {
		return nil, nil
	}

	return m, nil
}

func (s *mosaicPartialServiceImpl) Insert(mosaicPartial *model.MosaicPartial) error {
	return s.DbMap().Insert(mosaicPartial)
}

func (s *mosaicPartialServiceImpl) GetRandomMissing(mosaic *model.Mosaic) *model.MacroPartial {
	sql := `
		select *
		from macro_partials ma
		where ma.macro_id = ?
		and not exists (
			select 1 from mosaic_partials mo
			where mo.mosaic_id = ?
			and mo.macro_id = ?
		)
		order by random()
		limit 1
	`
	var macroPartial model.MacroPartial
	err := s.DbMap().SelectOne(&macroPartial, sql, mosaic.MacroId, mosaic.Id, mosaic.MacroId)
	if err != nil {
		return nil
	}

	return macroPartial
}
