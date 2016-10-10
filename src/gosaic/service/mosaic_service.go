package service

import (
	"database/sql"
	"errors"
	"fmt"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type MosaicService interface {
	Service
	Get(int64) (*model.Mosaic, error)
	Insert(*model.Mosaic) error
	Update(*model.Mosaic) (int64, error)
	GetOneBy(string, ...interface{}) (*model.Mosaic, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Mosaic, error)
}

type mosaicServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMosaicService(dbMap *gorp.DbMap) MosaicService {
	return &mosaicServiceImpl{dbMap: dbMap}
}

func (s *mosaicServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.Mosaic{}, "mosaics").SetKeys(true, "id")
	return nil
}

func (s *mosaicServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *mosaicServiceImpl) Get(id int64) (*model.Mosaic, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.Mosaic{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	m, ok := c.(*model.Mosaic)
	if !ok {
		return nil, errors.New("Unable to type cast mosaic")
	}

	if m.Id == int64(0) {
		return nil, nil
	}

	return m, nil
}

func (s *mosaicServiceImpl) Insert(mosaic *model.Mosaic) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(mosaic)
}

func (s *mosaicServiceImpl) Update(mosaic *model.Mosaic) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(mosaic)
}

func (s *mosaicServiceImpl) GetOneBy(conditions string, params ...interface{}) (*model.Mosaic, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var mosaic model.Mosaic

	err := s.dbMap.SelectOne(&mosaic, fmt.Sprintf("select * from mosaics where %s limit 1", conditions), params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &mosaic, nil
}

func (s *mosaicServiceImpl) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	num, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from mosaics where %s limit 1", conditions), params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return num == 1, nil
}

func (s *mosaicServiceImpl) FindAll(order string) ([]*model.Mosaic, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from mosaics order by %s", order)

	var mosaics []*model.Mosaic
	_, err := s.dbMap.Select(&mosaics, sql, order)

	return mosaics, err
}
