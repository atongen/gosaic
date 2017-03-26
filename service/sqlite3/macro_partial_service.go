package sqlite3

import (
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"

	"gopkg.in/gorp.v1"
)

type macroPartialServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMacroPartialService(dbMap *gorp.DbMap) *macroPartialServiceSqlite3 {
	return &macroPartialServiceSqlite3{dbMap: dbMap}
}

func (s *macroPartialServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.MacroPartial{}, "macro_partials").SetKeys(true, "id")
	return nil
}

func (s *macroPartialServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *macroPartialServiceSqlite3) Insert(macroPartial *model.MacroPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	err := macroPartial.EncodePixels()
	if err != nil {
		return err
	}
	return s.dbMap.Insert(macroPartial)
}

func (s *macroPartialServiceSqlite3) Update(macroPartial *model.MacroPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	err := macroPartial.EncodePixels()
	if err != nil {
		return err
	}
	_, err = s.dbMap.Update(macroPartial)
	return err
}

func (s *macroPartialServiceSqlite3) Delete(macroPartial *model.MacroPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(macroPartial)
	return err
}

func (s *macroPartialServiceSqlite3) Get(id int64) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	macroPartial, err := s.dbMap.Get(model.MacroPartial{}, id)
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

func (s *macroPartialServiceSqlite3) GetOneBy(column string, value interface{}) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var macroPartial model.MacroPartial
	err := s.dbMap.SelectOne(&macroPartial, fmt.Sprintf("select * from macro_partials where %s = ? limit 1", column), value)
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

func (s *macroPartialServiceSqlite3) ExistsBy(column string, value interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from macro_partials where %s = ? limit 1", column), value)
	return count == 1, err
}

func (s *macroPartialServiceSqlite3) Count(macro *model.Macro) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from macro_partials where macro_id = ?", macro.Id)
}

func (s *macroPartialServiceSqlite3) FindAll(order string, limit, offset int, conditions string, params ...interface{}) ([]*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

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

func (s *macroPartialServiceSqlite3) doFind(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	p := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
	}

	err := s.dbMap.SelectOne(&p, "select * from macro_partials where macro_id = ? and cover_partial_id = ?", p.MacroId, p.CoverPartialId)
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

func (s *macroPartialServiceSqlite3) Find(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(macro, coverPartial)
}

func (s *macroPartialServiceSqlite3) doCreate(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
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

	err = s.dbMap.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *macroPartialServiceSqlite3) Create(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(macro, coverPartial)
}

func (s *macroPartialServiceSqlite3) FindOrCreate(macro *model.Macro, coverPartial *model.CoverPartial) (*model.MacroPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(macro, coverPartial)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(macro, coverPartial)
}

func (s *macroPartialServiceSqlite3) CountMissing(macro *model.Macro) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
		select count(*)
		from cover_partials
		where cover_partials.cover_id = ?
		and not exists (
			select 1 from macro_partials
			where macro_partials.macro_id = ?
			and macro_partials.cover_partial_id = cover_partials.id
		)
	`

	return s.dbMap.SelectInt(sql, macro.CoverId, macro.Id)
}

func (s *macroPartialServiceSqlite3) FindMissing(macro *model.Macro, order string, limit, offset int) ([]*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`
		select *
		from cover_partials
		where cover_partials.cover_id = ?
		and not exists (
			select 1 from macro_partials
			where macro_partials.macro_id = ?
			and macro_partials.cover_partial_id = cover_partials.id
		)
		order by %s
		limit %d
		offset %d
	`, order, limit, offset)

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, macro.CoverId, macro.Id)

	return coverPartials, err
}

func (s *macroPartialServiceSqlite3) AspectIds(macroId int64) ([]int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

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
