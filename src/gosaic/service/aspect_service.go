package service

import (
	"bytes"
	"strconv"
	"sync"

	"gopkg.in/gorp.v1"

	"gosaic/model"
)

type AspectService interface {
	Service
	Insert(*model.Aspect) error
	Get(int64) (*model.Aspect, error)
	Count() (int64, error)
	Find(int, int) (*model.Aspect, error)
	Create(int, int) (*model.Aspect, error)
	FindOrCreate(int, int) (*model.Aspect, error)
	FindIn([]int64) ([]*model.Aspect, error)
}

type aspectServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewAspectService(dbMap *gorp.DbMap) AspectService {
	return &aspectServiceImpl{dbMap: dbMap}
}

func (s *aspectServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.Aspect{}, "aspects").SetKeys(true, "id")
	return nil
}

func (s *aspectServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *aspectServiceImpl) Insert(aspect *model.Aspect) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(aspect)
}

func (s *aspectServiceImpl) Get(id int64) (*model.Aspect, error) {
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

func (s *aspectServiceImpl) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from aspects")
}

func (s *aspectServiceImpl) Find(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(width, height)
}

func (s *aspectServiceImpl) doFind(width int, height int) (*model.Aspect, error) {
	aspect := model.NewAspect(width, height)

	err := s.dbMap.SelectOne(aspect, "select * from aspects where columns = ? and rows = ?", aspect.Columns, aspect.Rows)
	if err == nil {
		return aspect, nil
	}

	return nil, nil
}

func (s *aspectServiceImpl) Create(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(width, height)
}

func (s *aspectServiceImpl) doCreate(width int, height int) (*model.Aspect, error) {
	aspect := model.NewAspect(width, height)

	err := s.dbMap.Insert(aspect)
	if err != nil {
		return nil, err
	}

	return aspect, nil
}

func (s *aspectServiceImpl) FindOrCreate(width int, height int) (*model.Aspect, error) {
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

func (s *aspectServiceImpl) FindIn(ids []int64) ([]*model.Aspect, error) {
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
