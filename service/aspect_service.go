package service

import (
	"sync"

	"github.com/coopernurse/gorp"

	"github.com/atongen/gosaic/model"
)

type AspectService interface {
	Service
	Insert(...*model.Aspect) error
	Get(id int64) (*model.Aspect, error)
	Count() (int64, error)
	FindOrCreate(rows int, columns int) (*model.Aspect, error)
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

func (s *aspectServiceImpl) Register() {
	s.DbMap().AddTableWithName(model.Aspect{}, "aspects").SetKeys(true, "id")
}

func (s *aspectServiceImpl) Insert(aspects ...*model.Aspect) error {
	return s.DbMap().Insert(model.AspectsToInterface(aspects)...)
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
	} else {
		return &aspect, nil
	}
}
