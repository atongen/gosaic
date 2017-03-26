package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type coverServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewCoverService(dbMap *gorp.DbMap) *coverServiceSqlite3 {
	return &coverServiceSqlite3{dbMap: dbMap}
}

func (s *coverServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Cover{}, "covers").SetKeys(true, "id")
	return nil
}

func (s *coverServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *coverServiceSqlite3) Get(id int64) (*model.Cover, error) {
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

func (s *coverServiceSqlite3) Insert(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(c)
}

func (s *coverServiceSqlite3) Update(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(c)
	return err
}

func (s *coverServiceSqlite3) Delete(c *model.Cover) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(c)
	return err
}

func (s *coverServiceSqlite3) GetOneBy(conditions string, params ...interface{}) (*model.Cover, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := fmt.Sprintf("select * from covers where %s limit 1", conditions)

	var cover model.Cover
	err := s.dbMap.SelectOne(&cover, sqlStr, params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &cover, err
}

func (s *coverServiceSqlite3) FindAll(order string) ([]*model.Cover, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from covers order by %s", order)

	var covers []*model.Cover
	_, err := s.dbMap.Select(&covers, sql)

	return covers, err
}
