package sqlite3

import (
	"bytes"
	"strconv"
	"sync"

	"gopkg.in/gorp.v1"

	"github.com/atongen/gosaic/model"
)

type aspectServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewAspectService(dbMap *gorp.DbMap) *aspectServiceSqlite3 {
	return &aspectServiceSqlite3{dbMap: dbMap}
}

func (s *aspectServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Aspect{}, "aspects").SetKeys(true, "id")
	return nil
}

func (s *aspectServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *aspectServiceSqlite3) Insert(aspect *model.Aspect) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(aspect)
}

func (s *aspectServiceSqlite3) Get(id int64) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	aspect, err := s.dbMap.Get(model.Aspect{}, id)
	if err != nil {
		return nil, err
	} else if aspect != nil {
		return aspect.(*model.Aspect), nil
	} else {
		return nil, nil
	}
}

func (s *aspectServiceSqlite3) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from aspects")
}

func (s *aspectServiceSqlite3) Find(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(width, height)
}

func (s *aspectServiceSqlite3) doFind(width int, height int) (*model.Aspect, error) {
	aspect := model.NewAspect(width, height)

	err := s.dbMap.SelectOne(aspect, "select * from aspects where columns = ? and rows = ? limit 1", aspect.Columns, aspect.Rows)
	if err == nil {
		return aspect, nil
	}

	return nil, nil
}

func (s *aspectServiceSqlite3) Create(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(width, height)
}

func (s *aspectServiceSqlite3) doCreate(width int, height int) (*model.Aspect, error) {
	aspect := model.NewAspect(width, height)

	err := s.dbMap.Insert(aspect)
	if err != nil {
		return nil, err
	}

	return aspect, nil
}

func (s *aspectServiceSqlite3) FindOrCreate(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	aspect, err := s.doFind(width, height)
	if err != nil {
		return nil, err
	} else if aspect != nil {
		return aspect, nil
	}

	return s.doCreate(width, height)
}

func (s *aspectServiceSqlite3) FindIn(ids []int64) ([]*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	aspects := make([]*model.Aspect, 0)
	num := len(ids)
	if num == 0 {
		return aspects, nil
	}

	var b bytes.Buffer
	b.WriteString("select * from aspects where id in (")
	idsStr := make([]interface{}, num)
	for i, id := range ids {
		idsStr[i] = strconv.FormatInt(id, 10)
		b.WriteString("?")
		if i < num-1 {
			b.WriteString(",")
		}
	}
	b.WriteString(")")

	_, err := s.dbMap.Select(&aspects, b.String(), idsStr...)
	if err != nil {
		return nil, err
	}

	return aspects, nil
}
