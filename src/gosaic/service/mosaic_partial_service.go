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
	CountMissing(*model.Mosaic) (int64, error)
	GetMissing(*model.Mosaic) *model.MacroPartial
	GetRandomMissing(*model.Mosaic) *model.MacroPartial
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

func (s *mosaicPartialServiceImpl) CountMissing(mosaic *model.Mosaic) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
		select count(*)
		from macro_partials map
		where map.macro_id = ?
		and not exists (
			select 1 from mosaic_partials mop
			where mop.mosaic_id = ?
			and mop.macro_partial_id = map.id
		)
	`
	return s.DbMap().SelectInt(sql, mosaic.MacroId, mosaic.Id)
}

func (s *mosaicPartialServiceImpl) GetMissing(mosaic *model.Mosaic) *model.MacroPartial {
	sql := `
		select *
		from macro_partials map
		where map.macro_id = ?
		and not exists (
			select 1 from mosaic_partials mop
			where mop.mosaic_id = ?
			and mop.macro_partial_id = map.id
		)
		order by map.id asc
		limit 1
	`
	var macroPartial model.MacroPartial
	err := s.DbMap().SelectOne(&macroPartial, sql, mosaic.MacroId, mosaic.Id)
	if err != nil {
		return nil
	}

	return macroPartial
}

func (s *mosaicPartialServiceImpl) GetRandomMissing(mosaic *model.Mosaic) *model.MacroPartial {
	sql := `
		select *
		from macro_partials map
		where map.macro_id = ?
		and not exists (
			select 1 from mosaic_partials mop
			where mop.mosaic_id = ?
			and mop.macro_partial_id = map.id
		)
		order by random()
		limit 1
	`
	var macroPartial model.MacroPartial
	err := s.DbMap().SelectOne(&macroPartial, sql, mosaic.MacroId, mosaic.Id)
	if err != nil {
		return nil
	}

	return macroPartial
}
