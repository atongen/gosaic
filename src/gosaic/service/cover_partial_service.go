package service

import (
	"fmt"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type CoverPartialService interface {
	Service
	Get(int64) (*model.CoverPartial, error)
	Insert(*model.CoverPartial) error
	Update(*model.CoverPartial) error
	Delete(*model.CoverPartial) error
	FindAll(int64, string) ([]*model.CoverPartial, error)
}

type coverPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewCoverPartialService(dbMap *gorp.DbMap) CoverPartialService {
	return &coverPartialServiceImpl{dbMap: dbMap}
}

func (s *coverPartialServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.CoverPartial{}, "cover_partials").SetKeys(true, "id")
	return nil
}

func (s *coverPartialServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *coverPartialServiceImpl) Get(id int64) (*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.CoverPartial{}, id)
	if err != nil {
		return nil, err
	} else if c != nil {
		return c.(*model.CoverPartial), nil
	} else {
		return nil, nil
	}
}

func (s *coverPartialServiceImpl) Insert(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(c)
}

func (s *coverPartialServiceImpl) Update(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(c)
	return err
}

func (s *coverPartialServiceImpl) Delete(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(c)
	return err
}

func (s *coverPartialServiceImpl) FindAll(coverId int64, order string) ([]*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from cover_partials where cover_id = ? order by %s", order)

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, coverId)

	return coverPartials, err
}
