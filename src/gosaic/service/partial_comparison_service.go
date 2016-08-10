package service

import (
	"fmt"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type PartialComparisonService interface {
	Service
	Insert(*model.PartialComparison) error
	Update(*model.PartialComparison) error
	Delete(*model.PartialComparison) error
	Get(int64) (*model.PartialComparison, error)
	GetOneBy(string, ...interface{}) (*model.PartialComparison, error)
	ExistsBy(string, ...interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, ...interface{}) (int64, error)
	FindAll(string, int, int, string, ...interface{}) ([]*model.PartialComparison, error)
	Find(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	Create(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	FindOrCreate(*model.MacroPartial, *model.GidxPartial) (*model.PartialComparison, error)
	CountMissing(macro *model.Macro) (int64, error)
	FindMissing(*model.Macro, int) ([]*model.MacroGidxView, error)
}

type partialComparisonServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewPartialComparisonService(dbMap *gorp.DbMap) PartialComparisonService {
	return &partialComparisonServiceImpl{dbMap: dbMap}
}

func (s *partialComparisonServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *partialComparisonServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.PartialComparison{}, "partial_comparisons").SetKeys(true, "id")
	return nil
}

func (s *partialComparisonServiceImpl) Insert(pc *model.PartialComparison) error {
	return s.DbMap().Insert(pc)
}

func (s *partialComparisonServiceImpl) Update(pc *model.PartialComparison) error {
	_, err := s.DbMap().Update(pc)
	return err
}

func (s *partialComparisonServiceImpl) Delete(pc *model.PartialComparison) error {
	_, err := s.DbMap().Delete(pc)
	return err
}

func (s *partialComparisonServiceImpl) Get(id int64) (*model.PartialComparison, error) {
	partialComparison, err := s.DbMap().Get(model.PartialComparison{}, id)
	if err != nil {
		return nil, err
	}

	if partialComparison == nil {
		return nil, nil
	}

	mp, ok := partialComparison.(*model.PartialComparison)
	if !ok {
		return nil, fmt.Errorf("Received struct is not a PartialComparison")
	}

	if mp.Id == int64(0) {
		return nil, nil
	}

	return mp, nil
}

func (s *partialComparisonServiceImpl) GetOneBy(conditions string, params ...interface{}) (*model.PartialComparison, error) {
	var partialComparison model.PartialComparison
	err := s.DbMap().SelectOne(&partialComparison, fmt.Sprintf("select * from partial_comparisons where %s limit 1", conditions), params...)
	if err != nil {
		return nil, err
	}

	if partialComparison.Id == int64(0) {
		return nil, nil
	}

	return &partialComparison, nil
}

func (s *partialComparisonServiceImpl) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt(fmt.Sprintf("select 1 from partial_comparisons where %s limit 1", conditions), params)
	return count == 1, err
}

func (s *partialComparisonServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from partial_comparisons")
}

func (s *partialComparisonServiceImpl) CountBy(conditions string, params ...interface{}) (int64, error) {
	return s.DbMap().SelectInt(fmt.Sprintf("select count(*) from partial_comparisons where %s", conditions), params)
}

func (s *partialComparisonServiceImpl) FindAll(order string, limit, offset int, conditions string, params ...interface{}) ([]*model.PartialComparison, error) {
	var partialComparisons []*model.PartialComparison

	sql := fmt.Sprintf("select * from partial_comparisons where %s order by %s limit %d offset %d",
		conditions, order, limit, offset)

	_, err := s.dbMap.Select(&partialComparisons, sql, params...)
	if err != nil {
		return nil, err
	}

	return partialComparisons, nil
}

func (s *partialComparisonServiceImpl) doFind(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	p := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err := s.DbMap().SelectOne(&p, "select * from partial_comparisons where macro_partial_id = ? and gidx_partial_id = ? limit 1", p.MacroPartialId, p.GidxPartialId)
	if err != nil {
		return nil, err
	}

	if p.Id == int64(0) {
		return nil, nil
	}

	return &p, nil
}

func (s *partialComparisonServiceImpl) Find(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceImpl) doCreate(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	p := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	dist, err := model.PixelDist(macroPartial, gidxPartial)
	if err != nil {
		return nil, err
	}
	p.Dist = dist

	err = s.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *partialComparisonServiceImpl) Create(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceImpl) FindOrCreate(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(macroPartial, gidxPartial)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceImpl) CountMissing(macro *model.Macro) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
select count(*)
from macro_partials, gidx_partials
where macro_partials.macro_id = ?
and macro_partials.aspect_id = gidx_partials.aspect_id
and not exists (
	select 1 from partial_comparisons
	where partial_comparisons.macro_partial_id = macro_partials.id
	and partial_comparisons.gidx_partial_id = gidx_partials.id
)
`

	return s.DbMap().SelectInt(sql, macro.Id)
}

func (s *partialComparisonServiceImpl) FindMissing(macro *model.Macro, limit int) ([]*model.MacroGidxView, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
select macro_partials.id as macro_partial_id,
	macro_partials.macro_id,
	macro_partials.cover_partial_id,
	macro_partials.aspect_id,
	macro_partials.data as macro_partial_data,
	gidx_partials.id as gidx_partial_id,
	gidx_partials.gidx_id,
	gidx_partials.data as gidx_partial_data
from macro_partials join gidx_partials
where macro_partials.macro_id = ?
and macro_partials.aspect_id = gidx_partials.aspect_id
and not exists (
	select 1 from partial_comparisons
	where partial_comparisons.macro_partial_id = macro_partials.id
	and partial_comparisons.gidx_partial_id = gidx_partials.id
)
order by macro_partials.id asc,
	gidx_partials.id asc
limit ?
`

	var macroGidxViews []*model.MacroGidxView
	_, err := s.dbMap.Select(&macroGidxViews, sql, macro.Id, limit)

	return macroGidxViews, err
}
