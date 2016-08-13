package service

import (
	"strconv"
	"strings"
	"sync"

	"gopkg.in/gorp.v1"

	"gosaic/model"
)

type AspectService interface {
	Service
	Insert(*model.Aspect) error
	Get(int64) (*model.Aspect, error)
	Count() (int64, error)
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

func (s *aspectServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *aspectServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Aspect{}, "aspects").SetKeys(true, "id")
	return nil
}

func (s *aspectServiceImpl) Insert(aspect *model.Aspect) error {
	return s.DbMap().Insert(aspect)
}

func (s *aspectServiceImpl) Get(id int64) (*model.Aspect, error) {
	aspect, err := s.DbMap().Get(model.Aspect{}, id)
	if err != nil {
		return nil, err
	} else if aspect != nil {
		return aspect.(*model.Aspect), nil
	} else {
		return nil, nil
	}
}

func (s *aspectServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from aspects")
}

func (s *aspectServiceImpl) FindOrCreate(width int, height int) (*model.Aspect, error) {
	s.m.Lock()
	defer s.m.Unlock()

	aspect := model.Aspect{}
	aspect.SetAspect(width, height)

	// find
	err := s.DbMap().SelectOne(&aspect, "select * from aspects where columns = ? and rows = ?", aspect.Columns, aspect.Rows)
	if err == nil {
		return &aspect, nil
	}

	// or create
	err = s.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	return &aspect, nil
}

func (s *aspectServiceImpl) FindIn(ids []int64) ([]*model.Aspect, error) {
	aspects := make([]*model.Aspect, len(ids))
	if len(aspects) == 0 {
		return aspects, nil
	}

	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.FormatInt(id, 10)
	}
	_, err := s.DbMap().Select(&aspects, "select * from aspects where id in (?)", strings.Join(idsStr, ","))
	if err != nil {
		return nil, err
	}

	return aspects, nil
}
