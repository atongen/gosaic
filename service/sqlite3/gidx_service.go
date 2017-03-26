package sqlite3

import (
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"
	"gopkg.in/gorp.v1"
)

type gidxServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewGidxService(dbMap *gorp.DbMap) *gidxServiceSqlite3 {
	return &gidxServiceSqlite3{dbMap: dbMap}
}

func (s *gidxServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Gidx{}, "gidx").SetKeys(true, "id")
	return nil
}

func (s *gidxServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *gidxServiceSqlite3) Insert(gidx *model.Gidx) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(gidx)
}

func (s *gidxServiceSqlite3) Update(gidx *model.Gidx) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(gidx)
}

func (s *gidxServiceSqlite3) Delete(gidx *model.Gidx) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Delete(gidx)
}

func (s *gidxServiceSqlite3) Get(id int64) (*model.Gidx, error) {
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

func (s *gidxServiceSqlite3) GetOneBy(column string, value interface{}) (*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var gidx model.Gidx
	err := s.dbMap.SelectOne(&gidx, fmt.Sprintf("select * from gidx where %s = ? limit 1", column), value)
	return &gidx, err
}

func (s *gidxServiceSqlite3) ExistsBy(column string, value interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from gidx where %s = ?", column), value)
	return count == 1, err
}

func (s *gidxServiceSqlite3) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from gidx")
}

func (s *gidxServiceSqlite3) CountBy(column string, value interface{}) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt(fmt.Sprintf("select count(*) from gidx where %s = ?", column), value)
}

func (s *gidxServiceSqlite3) FindAll(order string, limit, offset int) ([]*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from gidx order by %s limit %d offset %d", order, limit, offset)

	var gidxs []*model.Gidx
	_, err := s.dbMap.Select(&gidxs, sql)

	return gidxs, err
}
