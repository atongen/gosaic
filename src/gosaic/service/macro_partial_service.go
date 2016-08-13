package service

import (
	"fmt"
	"gosaic/model"
	"gosaic/util"
	"sync"

	"gopkg.in/gorp.v1"
)

type MacroPartialService interface {
	Service
	Insert(*model.MacroPartial) error
	Update(*model.MacroPartial) error
	Delete(*model.MacroPartial) error
	Get(int64) (*model.MacroPartial, error)
	GetOneBy(string, interface{}) (*model.MacroPartial, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
	FindAll(string, int, int, string, ...interface{}) ([]*model.MacroPartial, error)
	Find(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	Create(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	FindOrCreate(*model.Macro, *model.CoverPartial) (*model.MacroPartial, error)
	FindMissing(*model.Macro, string, int, int) ([]*model.CoverPartial, error)
	AspectIds(int64) ([]int64, error)
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
	err := macroPartial.EncodePixels()
	if err != nil {
		return err
	}
	return s.DbMap().Insert(macroPartial)
}

func (s *macroPartialServiceImpl) Update(macroPartial *model.MacroPartial) error {
	err := macroPartial.EncodePixels()
	if err != nil {
		return err
	}
	_, err = s.DbMap().Update(macroPartial)
	return err
}

func (s *macroPartialServiceImpl) Delete(macroPartial *model.MacroPartial) error {
	_, err := s.DbMap().Delete(macroPartial)
	return err
}

func (s *macroPartialServiceImpl) Get(id int64) (*model.MacroPartial, error) {
	macroPartial, err := s.DbMap().Get(model.MacroPartial{}, id)
	if err != nil {
		return nil, err
	}

	if macroPartial == nil {
		return nil, nil
	}

	mp, ok := macroPartial.(*model.MacroPartial)
	if !ok {
		return nil, fmt.Errorf("Received struct is not a MacroPartial")
	}

	if mp.Id == int64(0) {
		return nil, nil
	}

	err = mp.DecodeData()
	if err != nil {
		return nil, err
	}

	return mp, nil
}

func (s *macroPartialServiceImpl) GetOneBy(column string, value interface{}) (*model.MacroPartial, error) {
	var macroPartial model.MacroPartial
	err := s.DbMap().SelectOne(&macroPartial, fmt.Sprintf("select * from macro_partials where %s = ? limit 1", column), value)
	if err != nil {
		return nil, err
	}

	if macroPartial.Id == int64(0) {
		return nil, nil
	}

	err = macroPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	return &macroPartial, err
}

func (s *macroPartialServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt(fmt.Sprintf("select 1 from macro_partials where %s = ? limit 1", column), value)
	return count == 1, err
}

func (s *macroPartialServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from macro_partials")
}

func (s *macroPartialServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt(fmt.Sprintf("select count(*) from macro_partials where %s = ?", column), value)
}

func (s *macroPartialServiceImpl) FindAll(order string, limit, offset int, conditions string, params ...interface{}) ([]*model.MacroPartial, error) {
	var macroPartials []*model.MacroPartial

	sql := fmt.Sprintf("select * from macro_partials where %s order by %s limit %d offset %d",
		conditions, order, limit, offset)

	_, err := s.dbMap.Select(&macroPartials, sql, params...)
	if err != nil {
		return nil, err
	}

	for _, mp := range macroPartials {
		err = mp.DecodeData()
		if err != nil {
			return nil, err
		}
	}

	return macroPartials, nil
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

	if p.Id == int64(0) {
		return nil, nil
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

func (s *macroPartialServiceImpl) doCreate(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	p := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       coverPartial.AspectId,
	}

	pixels, err := util.GetPartialLab(macro, coverPartial)
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

func (s *macroPartialServiceImpl) Create(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(macro, coverPartial)
}

func (s *macroPartialServiceImpl) FindOrCreate(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(macro, coverPartial)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(macro, coverPartial)
}

func (s *macroPartialServiceImpl) FindMissing(macro *model.Macro, order string, limit, offset int) ([]*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
select * from cover_partials
where cover_partials.cover_id = ?
and not exists (
	select 1 from macro_partials
	where macro_partials.macro_id = ?
	and macro_partials.cover_partial_id = cover_partials.id
)
order by ?
limit ?
offset ?
`

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, macro.CoverId, macro.Id, order, limit, offset)

	return coverPartials, err
}

func (s *macroPartialServiceImpl) AspectIds(macroId int64) ([]int64, error) {
	sql := `
		select distinct aspect_id
		from macro_partials
		where macro_id = ?
		order by aspect_id asc
	`
	aspectIds := make([]int64, 0)
	rows, err := s.dbMap.Db.Query(sql, macroId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var aspectId int64
		err = rows.Scan(&aspectId)
		if err != nil {
			return nil, err
		}
		aspectIds = append(aspectIds, aspectId)
	}

	return aspectIds, nil
}
