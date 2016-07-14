package service

import (
	"gosaic/model"
	"gosaic/util"
	"sync"

	"gopkg.in/gorp.v1"
)

type GidxPartialService interface {
	Service
	Insert(...*model.GidxPartial) error
	Update(...*model.GidxPartial) (int64, error)
	Delete(...*model.GidxPartial) (int64, error)
	Get(int64) (*model.GidxPartial, error)
	GetOneBy(string, interface{}) (*model.GidxPartial, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	FindOrCreate(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
}

type gidxPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewGidxPartialService(dbMap *gorp.DbMap) GidxPartialService {
	return &gidxPartialServiceImpl{dbMap: dbMap}
}

func (s *gidxPartialServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *gidxPartialServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.GidxPartial{}, "gidx_partials").SetKeys(true, "id")
	return nil
}

func (s *gidxPartialServiceImpl) Insert(gidx_partials ...*model.GidxPartial) error {
	return s.DbMap().Insert(model.GidxPartialsToInterface(gidx_partials)...)
}

func (s *gidxPartialServiceImpl) Update(gidx_partials ...*model.GidxPartial) (int64, error) {
	return s.DbMap().Update(model.GidxPartialsToInterface(gidx_partials)...)
}

func (s *gidxPartialServiceImpl) Delete(gidx_partials ...*model.GidxPartial) (int64, error) {
	return s.DbMap().Delete(model.GidxPartialsToInterface(gidx_partials)...)
}

func (s *gidxPartialServiceImpl) Get(id int64) (*model.GidxPartial, error) {
	gidx_partial, err := s.DbMap().Get(model.GidxPartial{}, id)
	if err != nil {
		return nil, err
	} else if gidx_partial != nil {
		return gidx_partial.(*model.GidxPartial), nil
	} else {
		return nil, nil
	}
}

func (s *gidxPartialServiceImpl) GetOneBy(column string, value interface{}) (*model.GidxPartial, error) {
	var gidx_partial model.GidxPartial
	err := s.DbMap().SelectOne(&gidx_partial, "select * from gidx_partials where "+column+" = ?", value)
	return &gidx_partial, err
}

func (s *gidxPartialServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from gidx_partials where "+column+" = ?", value)
	return count == 1, err
}

func (s *gidxPartialServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx_partials")
}

func (s *gidxPartialServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx_partials where "+column+" = ?", value)
}

func (s *gidxPartialServiceImpl) FindOrCreate(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	// find
	err := s.DbMap().SelectOne(&p, "select * from gidx_partials where gidx_id = ? and aspect_id = ?", p.GidxId, p.AspectId)
	if err == nil {
		err = p.DecodeData()
		if err != nil {
			return nil, err
		}
		return &p, nil
	}

	pixels, err := util.GetAspectLab(gidx, aspect)
	if err != nil {
		return nil, err
	}
	p.Pixels = pixels

	err = p.EncodePixels()
	if err != nil {
		return nil, err
	}

	// or create
	err = s.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
