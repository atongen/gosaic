package service

import (
	"gosaic/model"
	"gosaic/util"
	"sync"

	"gopkg.in/gorp.v1"
)

type MacroPartialService interface {
	Service
	Insert(*model.MacroPartial) error
	Update(*model.MacroPartial) (int64, error)
	Delete(*model.MacroPartial) (int64, error)
	Get(int64) (*model.MacroPartial, error)
	GetOneBy(string, interface{}) (*model.MacroPartial, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	Find(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	Create(*model.Macro, *model.CoverPartial, *model.Aspect) (*model.MacroPartial, error)
	FindOrCreate(*model.Macro, *model.CoverPartial, *model.Aspect) (*model.MacroPartial, error)
	FindMissing(*model.Macro, string, int, int) ([]*model.CoverPartial, error)
}

type macroPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMacroPartialService(dbMap *gorp.DbMap) MacroPartialService {
	return &macroPartialServiceImpl{dbMap: dbMap}
}

func (s *macroPartialServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *macroPartialServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.MacroPartial{}, "macro_partials").SetKeys(true, "id")
	return nil
}

func (s *macroPartialServiceImpl) Insert(macroPartial *model.MacroPartial) error {
	return s.DbMap().Insert(macroPartial)
}

func (s *macroPartialServiceImpl) Update(macroPartial *model.MacroPartial) (int64, error) {
	return s.DbMap().Update(macroPartial)
}

func (s *macroPartialServiceImpl) Delete(macroPartial *model.MacroPartial) (int64, error) {
	return s.DbMap().Delete(macroPartial)
}

func (s *macroPartialServiceImpl) Get(id int64) (*model.MacroPartial, error) {
	macro_partial, err := s.DbMap().Get(model.MacroPartial{}, id)
	if err != nil {
		return nil, err
	} else if macro_partial != nil {
		return macro_partial.(*model.MacroPartial), nil
	} else {
		return nil, nil
	}
}

func (s *macroPartialServiceImpl) GetOneBy(column string, value interface{}) (*model.MacroPartial, error) {
	var macro_partial model.MacroPartial
	err := s.DbMap().SelectOne(&macro_partial, "select * from macro_partials where "+column+" = ?", value)
	return &macro_partial, err
}

func (s *macroPartialServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from macro_partials where "+column+" = ?", value)
	return count == 1, err
}

func (s *macroPartialServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from macro_partials")
}

func (s *macroPartialServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from macro_partials where "+column+" = ?", value)
}

func (s *macroPartialServiceImpl) doFind(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	p := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
	}

	err := s.DbMap().SelectOne(&p, "select * from macro_partials where macro_id = ? and cover_partial_id = ?", p.MacroId, p.CoverPartialId)
	if err != nil {
		return nil, err
	}

	err = p.DecodeData()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *macroPartialServiceImpl) Find(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(macro, coverPartial)
}

func (s *macroPartialServiceImpl) doCreate(macro *model.Macro, coverPartial *model.CoverPartial, aspect *model.Aspect) (*model.MacroPartial, error) {
	p := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       coverPartial.AspectId,
	}

	pixels, err := util.GetAspectLab(macro, aspect)
	if err != nil {
		return nil, err
	}
	p.Pixels = pixels

	err = p.EncodePixels()
	if err != nil {
		return nil, err
	}

	err = s.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *macroPartialServiceImpl) Create(macro *model.Macro, coverPartial *model.CoverPartial, aspect *model.Aspect) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(macro, coverPartial, aspect)
}

func (s *macroPartialServiceImpl) FindOrCreate(macro *model.Macro, coverPartial *model.CoverPartial, aspect *model.Aspect) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(macro, coverPartial)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(macro, coverPartial, aspect)
}

func (s *macroPartialServiceImpl) FindMissing(macro *model.Macro, order string, limit, offset int) ([]*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
select * from cover_partials
where not exists (
	select 1 from macro_partials
	where macro_partials.macro_id = ?
	and macro_partials.cover_partial_id = cover_partials.id
)
order by ?
limit ?
offset ?
`

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, macro.Id, order, limit, offset)

	return coverPartials, err
}
