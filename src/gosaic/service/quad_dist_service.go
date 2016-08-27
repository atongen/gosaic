package service

import (
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type QuadDistService interface {
	Service
	Get(int64) (*model.QuadDist, error)
	Insert(*model.QuadDist) error
	GetWorst(*model.Macro) (*model.CoverPartial, error)
}

type quadDistServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewQuadDistService(dbMap *gorp.DbMap) QuadDistService {
	return &quadDistServiceImpl{dbMap: dbMap}
}

func (s *quadDistServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.QuadDist{}, "quad_dists").SetKeys(true, "id")
	return nil
}

func (s *quadDistServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *quadDistServiceImpl) Get(id int64) (*model.QuadDist, error) {
	quadDist, err := s.dbMap.Get(model.QuadDist{}, id)
	if err != nil {
		return nil, err
	} else if quadDist != nil {
		return quadDist.(*model.QuadDist), nil
	} else {
		return nil, nil
	}
}

func (s *quadDistServiceImpl) Insert(pc *model.QuadDist) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(pc)
}

func (s *quadDistServiceImpl) GetWorst(macro *model.Macro) (*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
		select cop.*
		from quad_dists qd
		inner join macro_partials map
			on qd.macro_partial_id = map.id
		inner join cover_partials cop
			on map.cover_partial_id = cop.id
		where map.macro_id = ?
		order by qd.dist desc
		limit 1
	`
	var coverPartial model.CoverPartial
	err := s.dbMap.SelectOne(&coverPartial, sql, macro.Id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &coverPartial, nil
}
