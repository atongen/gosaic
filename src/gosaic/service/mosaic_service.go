package service

import (
	"errors"
	"fmt"
	"gosaic/model"
	"sync"
	"time"

	"gopkg.in/gorp.v1"
)

type MosaicService interface {
	Service
	Get(int64) (*model.Mosaic, error)
	Insert(*model.Mosaic) error
	GetOneBy(string, ...interface{}) (*model.Mosaic, error)
	FindAll(string) ([]*model.Mosaic, error)
}

type mosaicServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMosaicService(dbMap *gorp.DbMap) MosaicService {
	return &mosaicServiceImpl{dbMap: dbMap}
}

func (s *mosaicServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *mosaicServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Mosaic{}, "mosaics").SetKeys(true, "id")
	return nil
}

func (s *mosaicServiceImpl) Get(id int64) (*model.Mosaic, error) {
	c, err := s.DbMap().Get(model.Mosaic{}, id)
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
	if mosaic.CreatedAt.IsZero() {
		mosaic.CreatedAt = time.Now()
	}
	return s.DbMap().Insert(mosaic)
}

func (s *mosaicServiceImpl) GetOneBy(conditions string, params ...interface{}) (*model.Mosaic, error) {
	var mosaic model.Mosaic

	err := s.DbMap().SelectOne(&mosaic, fmt.Sprintf("select * from mosaics where %s limit 1", conditions), params...)
	// returns error if none are found
	// or if more than one is found
	if err != nil {
		return nil, nil
	}

	return &mosaic, nil
}

func (s *mosaicServiceImpl) FindAll(order string) ([]*model.Mosaic, error) {
	sql := `select * from mosaics order by ?`

	var mosaics []*model.Mosaic
	_, err := s.dbMap.Select(&mosaics, sql, order)

	return mosaics, err
}
