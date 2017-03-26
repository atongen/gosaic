package sqlite3

import (
	"bytes"
	"database/sql"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type partialComparisonServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewPartialComparisonService(dbMap *gorp.DbMap) *partialComparisonServiceSqlite3 {
	return &partialComparisonServiceSqlite3{dbMap: dbMap}
}

func (s *partialComparisonServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.PartialComparison{}, "partial_comparisons").SetKeys(true, "id")
	return nil
}

func (s *partialComparisonServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *partialComparisonServiceSqlite3) Insert(pc *model.PartialComparison) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(pc)
}

func (s *partialComparisonServiceSqlite3) BulkInsert(partialComparisons []*model.PartialComparison) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(partialComparisons) == 0 {
		return int64(0), nil
	} else if len(partialComparisons) == 1 {
		err := s.dbMap.Insert(partialComparisons[0])
		if err != nil {
			return int64(0), err
		}
		return int64(1), nil
	}

	var b bytes.Buffer

	b.WriteString("insert into partial_comparisons (id, macro_partial_id, gidx_partial_id, dist) ")
	b.WriteString(fmt.Sprintf("select null as id, %d as macro_partial_id, %d as gidx_partial_id, %f as dist",
		partialComparisons[0].MacroPartialId, partialComparisons[0].GidxPartialId, partialComparisons[0].Dist))

	for i := 1; i < len(partialComparisons); i++ {
		b.WriteString(fmt.Sprintf(" union select null, %d, %d, %f",
			partialComparisons[i].MacroPartialId, partialComparisons[i].GidxPartialId, partialComparisons[i].Dist))
	}

	res, err := s.dbMap.Db.Exec(b.String())
	if err != nil {
		return int64(0), err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return int64(0), err
	}

	return rowsAffected, nil
}

func (s *partialComparisonServiceSqlite3) Update(pc *model.PartialComparison) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(pc)
	return err
}

func (s *partialComparisonServiceSqlite3) Delete(pc *model.PartialComparison) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(pc)
	return err
}

func (s *partialComparisonServiceSqlite3) DeleteBy(conditions string, params ...interface{}) error {
	s.m.Lock()
	defer s.m.Unlock()

	if conditions == "" || len(params) == 0 {
		return nil
	}

	sqlStr := fmt.Sprintf("delete from partial_comparisons where %s;", conditions)
	_, err := s.dbMap.Db.Exec(sqlStr, params...)
	return err
}

func (s *partialComparisonServiceSqlite3) DeleteFrom(macro *model.Macro) error {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := `
		delete from partial_comparisons
		where id in (
			select pcs.id
			from partial_comparisons pcs
			inner join macro_partials maps
				on pcs.macro_partial_id = maps.id
			where maps.macro_id = ?
		)
	`
	_, err := s.dbMap.Db.Exec(sqlStr, macro.Id)
	return err
}

func (s *partialComparisonServiceSqlite3) Get(id int64) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	partialComparison, err := s.dbMap.Get(model.PartialComparison{}, id)
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

func (s *partialComparisonServiceSqlite3) GetOneBy(conditions string, params ...interface{}) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var partialComparison model.PartialComparison
	err := s.dbMap.SelectOne(&partialComparison, fmt.Sprintf("select * from partial_comparisons where %s limit 1", conditions), params...)
	if err != nil {
		return nil, err
	}

	return &partialComparison, nil
}

func (s *partialComparisonServiceSqlite3) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from partial_comparisons where %s limit 1", conditions), params...)
	return count == 1, err
}

func (s *partialComparisonServiceSqlite3) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from partial_comparisons")
}

func (s *partialComparisonServiceSqlite3) CountBy(conditions string, params ...interface{}) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`
		select count(*)
		from partial_comparisons
		where %s
		limit 1
	`, conditions)

	return s.dbMap.SelectInt(sql, params...)
}

func (s *partialComparisonServiceSqlite3) FindAll(order string, limit, offset int, conditions string, params ...interface{}) ([]*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var partialComparisons []*model.PartialComparison

	sql := fmt.Sprintf("select * from partial_comparisons where %s order by %s limit %d offset %d",
		conditions, order, limit, offset)

	_, err := s.dbMap.Select(&partialComparisons, sql, params...)
	if err != nil {
		return nil, err
	}

	return partialComparisons, nil
}

func (s *partialComparisonServiceSqlite3) doFind(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	p := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err := s.dbMap.SelectOne(&p, "select * from partial_comparisons where macro_partial_id = ? and gidx_partial_id = ? limit 1", p.MacroPartialId, p.GidxPartialId)
	if err != nil {
		return nil, err
	}

	if p.Id == int64(0) {
		return nil, nil
	}

	return &p, nil
}

func (s *partialComparisonServiceSqlite3) Find(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceSqlite3) doCreate(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	p := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	dist, err := model.PixelDist(macroPartial, gidxPartial)
	if err != nil {
		return nil, err
	}
	p.Dist = dist

	err = s.dbMap.Insert(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *partialComparisonServiceSqlite3) Create(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceSqlite3) FindOrCreate(macroPartial *model.MacroPartial, gidxPartial *model.GidxPartial) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(macroPartial, gidxPartial)
	if err == nil {
		return p, nil
	}

	// or create
	return s.doCreate(macroPartial, gidxPartial)
}

func (s *partialComparisonServiceSqlite3) CountMissing(macro *model.Macro) (int64, error) {
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

	return s.dbMap.SelectInt(sql, macro.Id)
}

func (s *partialComparisonServiceSqlite3) FindMissing(macro *model.Macro, limit int) ([]*model.MacroGidxView, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`
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
limit %d
`, limit)

	var macroGidxViews []*model.MacroGidxView
	rows, err := s.dbMap.Db.Query(sql, macro.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.MacroGidxView = model.MacroGidxView{
			&model.MacroPartial{},
			&model.GidxPartial{},
		}
		err = rows.Scan(
			&r.MacroPartial.Id,
			&r.MacroPartial.MacroId,
			&r.MacroPartial.CoverPartialId,
			&r.MacroPartial.AspectId,
			&r.MacroPartial.Data,
			&r.GidxPartial.Id,
			&r.GidxPartial.GidxId,
			&r.GidxPartial.Data,
		)
		if err != nil {
			return nil, err
		}
		r.GidxPartial.AspectId = r.MacroPartial.AspectId
		macroGidxViews = append(macroGidxViews, &r)
	}

	return macroGidxViews, nil
}

func (s *partialComparisonServiceSqlite3) CreateFromView(view *model.MacroGidxView) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	pc, err := view.PartialComparison()
	if err != nil {
		return nil, err
	}

	err = s.dbMap.Insert(pc)
	if err != nil {
		return nil, err
	}

	return pc, nil
}

func (s *partialComparisonServiceSqlite3) GetClosest(macroPartial *model.MacroPartial) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := `
		select pc.gidx_partial_id
		from partial_comparisons pc
		where pc.macro_partial_id = ?
		order by pc.dist asc
		limit 1
	`
	gidxPartialId, err := s.dbMap.SelectInt(sqlStr, macroPartial.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return gidxPartialId, nil
}

func (s *partialComparisonServiceSqlite3) GetClosestMax(macroPartial *model.MacroPartial, mosaic *model.Mosaic, maxRepeats int) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := fmt.Sprintf(`
		select pc.gidx_partial_id
		from partial_comparisons pc
		where pc.macro_partial_id = ?
		and not exists (
			select 1
			from mosaic_partials mos
			inner join gidx_partials gps
			on mos.gidx_partial_id = gps.id
			where mos.gidx_partial_id = pc.gidx_partial_id
			and mos.mosaic_id = ?
			group by gps.gidx_id
			having count(*) >= %d
		)
		order by pc.dist asc
		limit 1
	`, maxRepeats)

	gidxPartialId, err := s.dbMap.SelectInt(sqlStr, macroPartial.Id, mosaic.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return gidxPartialId, nil
}

func (s *partialComparisonServiceSqlite3) GetBestAvailable(mosaic *model.Mosaic) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := `
		select pc.*
		from partial_comparisons pc
		inner join macro_partials map
			on pc.macro_partial_id = map.id
		where map.macro_id = ?
		and not exists (
			select 1
			from mosaic_partials mos
			where mos.mosaic_id = ?
			and mos.macro_partial_id = map.id
		)
		order by pc.dist asc
		limit 1
	`

	var partialComparison model.PartialComparison
	// returns error on no results
	err := s.dbMap.SelectOne(&partialComparison, sqlStr, mosaic.MacroId, mosaic.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &partialComparison, nil
}

func (s *partialComparisonServiceSqlite3) GetBestAvailableMax(mosaic *model.Mosaic, maxRepeats int) (*model.PartialComparison, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := fmt.Sprintf(`
		select pc.*
		from partial_comparisons pc
		join macro_partials map
		join mosaics mo
		where mo.id = ?
		and pc.macro_partial_id = map.id
		and map.macro_id = mo.macro_id
		and not exists (
			select 1
			from mosaic_partials mop
			where mop.mosaic_id = mo.id
			and mop.macro_partial_id = map.id
		) and not exists (
			select 1
			from mosaic_partials mop
			inner join gidx_partials gp
				on mop.gidx_partial_id = gp.id
			where mop.gidx_partial_id = pc.gidx_partial_id
			and mop.mosaic_id = mo.id
			group by gp.gidx_id
			having count(*) >= %d
		)
		order by pc.dist asc
		limit 1
	`, maxRepeats)

	var partialComparison model.PartialComparison
	// returns error on no results
	err := s.dbMap.SelectOne(&partialComparison, sqlStr, mosaic.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &partialComparison, nil
}
