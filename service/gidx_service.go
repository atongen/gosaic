package service

import (
	"github.com/coopernurse/gorp"

	"github.com/atongen/gosaic/model"
)

type GidxService interface {
	Insert(...*model.Gidx) error
	Update(...*model.Gidx) (int64, error)
	Delete(...*model.Gidx) (int64, error)
	Get(int64) (*model.Gidx, error)
	GetOneBy(string, interface{}) (*model.Gidx, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
}

type GidxServiceImpl struct {
	dbMap *gorp.DbMap
}

func NewGidxService(dbMap *gorp.DbMap) *GidxServiceImpl {
	return &GidxServiceImpl{dbMap: dbMap}
}

func (s *GidxServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *GidxServiceImpl) Register() {
	s.DbMap().AddTableWithName(model.Gidx{}, "gidx").SetKeys(true, "Id")
}

func (s *GidxServiceImpl) Insert(gidxs ...*model.Gidx) error {
	return s.DbMap().Insert(model.GidxsToInterface(gidxs)...)
}

func (s *GidxServiceImpl) Update(gidxs ...*model.Gidx) (int64, error) {
	return s.DbMap().Update(model.GidxsToInterface(gidxs)...)
}

func (s *GidxServiceImpl) Delete(gidxs ...*model.Gidx) (int64, error) {
	return s.DbMap().Delete(model.GidxsToInterface(gidxs)...)
}

func (s *GidxServiceImpl) Get(id int64) (*model.Gidx, error) {
	gidx, err := s.DbMap().Get(model.Gidx{}, id)
	if err != nil {
		return nil, err
	} else if gidx != nil {
		return gidx.(*model.Gidx), nil
	} else {
		return nil, nil
	}
}

func (s *GidxServiceImpl) GetOneBy(column string, value interface{}) (*model.Gidx, error) {
	var gidx model.Gidx
	err := s.DbMap().SelectOne(&gidx, "select * from \"gidx\" where "+column+" = ?", value)
	return &gidx, err
}

func (s *GidxServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from \"gidx\" where "+column+" = ?", value)
	return count == 1, err
}

func (s *GidxServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from \"gidx\"")
}

func (s *GidxServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from \"gidx\" where "+column+" = ?", value)
}
