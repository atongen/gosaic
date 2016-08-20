package service

import (
	"errors"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type CoverService interface {
	Service
	Get(int64) (*model.Cover, error)
	Insert(*model.Cover) error
	Update(*model.Cover) error
	Delete(*model.Cover) error
	GetOneBy(string, interface{}) (*model.Cover, error)
	FindAll(string) ([]*model.Cover, error)
}

type coverServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewCoverService(dbMap *gorp.DbMap) CoverService {
	return &coverServiceImpl{dbMap: dbMap}
}

func (s *coverServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.Cover{}, "covers").SetKeys(true, "id")
	return nil
}

func (s *coverServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *coverServiceImpl) Get(id int64) (*model.Cover, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.Cover{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	cc, ok := c.(*model.Cover)
	if !ok {
		return nil, errors.New("Unable to type cast cover")
	}

	if cc.Id == int64(0) {
		return nil, nil
	}

	return cc, nil
}

func (s *coverServiceImpl) Insert(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(c)
}

func (s *coverServiceImpl) Update(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(c)
	return err
}

func (s *coverServiceImpl) Delete(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(c)
	return err
}

func (s *coverServiceImpl) GetOneBy(column string, value interface{}) (*model.Cover, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var cover model.Cover
	err := s.dbMap.SelectOne(&cover, "select * from covers where "+column+" = ? limit 1", value)
	return &cover, err
}

func (s *coverServiceImpl) FindAll(order string) ([]*model.Cover, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `select * from covers order by ?`

	var covers []*model.Cover
	_, err := s.dbMap.Select(&covers, sql, order)

	return covers, err
}
