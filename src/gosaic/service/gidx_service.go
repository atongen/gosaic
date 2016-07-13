package service

import (
	"gosaic/model"

	"gopkg.in/gorp.v1"
)

type GidxService interface {
	Service
	Insert(...*model.Gidx) error
	Update(...*model.Gidx) (int64, error)
	Delete(...*model.Gidx) (int64, error)
	Get(int64) (*model.Gidx, error)
	GetOneBy(string, interface{}) (*model.Gidx, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
}

type gidxServiceImpl struct {
	dbMap *gorp.DbMap
}

func NewGidxService(dbMap *gorp.DbMap) GidxService {
	return &gidxServiceImpl{dbMap: dbMap}
}

func (s *gidxServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *gidxServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Gidx{}, "gidx").SetKeys(true, "id")
	return nil
}

func (s *gidxServiceImpl) Insert(gidxs ...*model.Gidx) error {
	return s.DbMap().Insert(model.GidxsToInterface(gidxs)...)
}

func (s *gidxServiceImpl) Update(gidxs ...*model.Gidx) (int64, error) {
	return s.DbMap().Update(model.GidxsToInterface(gidxs)...)
}

func (s *gidxServiceImpl) Delete(gidxs ...*model.Gidx) (int64, error) {
	return s.DbMap().Delete(model.GidxsToInterface(gidxs)...)
}

func (s *gidxServiceImpl) Get(id int64) (*model.Gidx, error) {
	gidx, err := s.DbMap().Get(model.Gidx{}, id)
	if err != nil {
		return nil, err
	} else if gidx != nil {
		return gidx.(*model.Gidx), nil
	} else {
		return nil, nil
	}
}

func (s *gidxServiceImpl) GetOneBy(column string, value interface{}) (*model.Gidx, error) {
	var gidx model.Gidx
	err := s.DbMap().SelectOne(&gidx, "select * from gidx where "+column+" = ?", value)
	return &gidx, err
}

func (s *gidxServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from gidx where "+column+" = ?", value)
	return count == 1, err
}

func (s *gidxServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx")
}

func (s *gidxServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from gidx where "+column+" = ?", value)
}
