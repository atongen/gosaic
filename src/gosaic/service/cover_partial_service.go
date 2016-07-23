package service

import (
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

func (s *coverPartialServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *coverPartialServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.CoverPartial{}, "cover_partials").SetKeys(true, "id")
	return nil
}

func (s *coverPartialServiceImpl) Get(id int64) (*model.CoverPartial, error) {
	c, err := s.DbMap().Get(model.CoverPartial{}, id)
	if err != nil {
		return nil, err
	} else if c != nil {
		return c.(*model.CoverPartial), nil
	} else {
		return nil, nil
	}
}

func (s *coverPartialServiceImpl) Insert(c *model.CoverPartial) error {
	return s.DbMap().Insert(c)
}

func (s *coverPartialServiceImpl) Update(c *model.CoverPartial) error {
	_, err := s.DbMap().Update(c)
	return err
}

func (s *coverPartialServiceImpl) Delete(c *model.CoverPartial) error {
	_, err := s.DbMap().Delete(c)
	return err
}

func (s *coverPartialServiceImpl) FindAll(coverId int64, order string) ([]*model.CoverPartial, error) {
	sql := `select * from cover_partials where cover_id = ? order by ?`

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, coverId, order)

	return coverPartials, err
}
