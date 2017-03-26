package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type macroServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMacroService(dbMap *gorp.DbMap) *macroServiceSqlite3 {
	return &macroServiceSqlite3{dbMap: dbMap}
}

func (s *macroServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Macro{}, "macros").SetKeys(true, "id")
	return nil
}

func (s *macroServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *macroServiceSqlite3) Get(id int64) (*model.Macro, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.Macro{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	m, ok := c.(*model.Macro)
	if !ok {
		return nil, errors.New("Unable to type cast macro")
	}

	if m.Id == int64(0) {
		return nil, nil
	}

	return m, nil
}

func (s *macroServiceSqlite3) Insert(c *model.Macro) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(c)
}

func (s *macroServiceSqlite3) Update(c *model.Macro) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(c)
	return err
}

func (s *macroServiceSqlite3) Delete(c *model.Macro) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(c)
	return err
}

func (s *macroServiceSqlite3) GetOneBy(conditions string, params ...interface{}) (*model.Macro, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var macro model.Macro

	err := s.dbMap.SelectOne(&macro, fmt.Sprintf("select * from macros where %s limit 1", conditions), params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &macro, nil
}

func (s *macroServiceSqlite3) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from macros where %s limit 1", conditions), params...)
	return count == 1, err
}

func (s *macroServiceSqlite3) FindAll(order string) ([]*model.Macro, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`select * from macros order by %s`, order)

	var macros []*model.Macro
	_, err := s.dbMap.Select(&macros, sql, order)

	return macros, err
}
