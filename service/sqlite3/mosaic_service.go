package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type mosaicServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMosaicService(dbMap *gorp.DbMap) *mosaicServiceSqlite3 {
	return &mosaicServiceSqlite3{dbMap: dbMap}
}

func (s *mosaicServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Mosaic{}, "mosaics").SetKeys(true, "id")
	return nil
}

func (s *mosaicServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *mosaicServiceSqlite3) Get(id int64) (*model.Mosaic, error) {
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

func (s *mosaicServiceSqlite3) Insert(mosaic *model.Mosaic) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(mosaic)
}

func (s *mosaicServiceSqlite3) Update(mosaic *model.Mosaic) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(mosaic)
}

func (s *mosaicServiceSqlite3) GetOneBy(conditions string, params ...interface{}) (*model.Mosaic, error) {
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

func (s *mosaicServiceSqlite3) ExistsBy(conditions string, params ...interface{}) (bool, error) {
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

func (s *mosaicServiceSqlite3) FindAll(order string) ([]*model.Mosaic, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from mosaics order by %s", order)

	var mosaics []*model.Mosaic
	_, err := s.dbMap.Select(&mosaics, sql, order)

	return mosaics, err
}
