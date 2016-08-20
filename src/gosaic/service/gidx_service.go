package service

import (
	"fmt"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type GidxService interface {
	Service
	Insert(*model.Gidx) error
	Update(*model.Gidx) (int64, error)
	Delete(*model.Gidx) (int64, error)
	Get(int64) (*model.Gidx, error)
	GetOneBy(string, interface{}) (*model.Gidx, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	FindAll(string, int, int) ([]*model.Gidx, error)
}

type gidxServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewGidxService(dbMap *gorp.DbMap) GidxService {
	return &gidxServiceImpl{dbMap: dbMap}
}

func (s *gidxServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.Gidx{}, "gidx").SetKeys(true, "id")
	return nil
}

func (s *gidxServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *gidxServiceImpl) Insert(gidx *model.Gidx) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(gidx)
}

func (s *gidxServiceImpl) Update(gidx *model.Gidx) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(gidx)
}

func (s *gidxServiceImpl) Delete(gidx *model.Gidx) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Delete(gidx)
}

func (s *gidxServiceImpl) Get(id int64) (*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	gidx, err := s.dbMap.Get(model.Gidx{}, id)
	if err != nil {
		return nil, err
	} else if gidx != nil {
		return gidx.(*model.Gidx), nil
	} else {
		return nil, nil
	}
}

func (s *gidxServiceImpl) GetOneBy(column string, value interface{}) (*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var gidx model.Gidx
	err := s.dbMap.SelectOne(&gidx, fmt.Sprintf("select * from gidx where %s = ? limit 1", column), value)
	return &gidx, err
}

func (s *gidxServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from gidx where %s = ?", column), value)
	return count == 1, err
}

func (s *gidxServiceImpl) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from gidx")
}

func (s *gidxServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt(fmt.Sprintf("select count(*) from gidx where %s = ?", column), value)
}

func (s *gidxServiceImpl) FindAll(order string, limit, offset int) ([]*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from gidx order by %s limit %d offset %d", order, limit, offset)

	var gidxs []*model.Gidx
	_, err := s.dbMap.Select(&gidxs, sql)

	return gidxs, err
}
