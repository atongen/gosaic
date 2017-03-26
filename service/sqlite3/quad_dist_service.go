package sqlite3

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type quadDistServiceSqlite struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewQuadDistService(dbMap *gorp.DbMap) *quadDistServiceSqlite {
	return &quadDistServiceSqlite{dbMap: dbMap}
}

func (s *quadDistServiceSqlite) Register() error {
	s.dbMap.AddTableWithName(model.QuadDist{}, "quad_dists").SetKeys(true, "id")
	return nil
}

func (s *quadDistServiceSqlite) Close() error {
	return s.dbMap.Db.Close()
}

func (s *quadDistServiceSqlite) Get(id int64) (*model.QuadDist, error) {
	s.m.Lock()
	defer s.m.Unlock()

	quadDist, err := s.dbMap.Get(model.QuadDist{}, id)
	if err != nil {
		return nil, err
	} else if quadDist != nil {
		return quadDist.(*model.QuadDist), nil
	} else {
		return nil, nil
	}
}

func (s *quadDistServiceSqlite) Insert(pc *model.QuadDist) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(pc)
}

func (s *quadDistServiceSqlite) GetWorst(macro *model.Macro, depth, area int) (*model.CoverPartialQuadView, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sqlStr := fmt.Sprintf(`
		select cop.id as cover_partial_id,
			cop.cover_id as cover_partial_cover_id,
			cop.aspect_id as cover_partial_aspect_id,
			cop.x1 as cover_partial_x1,
			cop.y1 as cover_partial_y1,
			cop.x2 as cover_partial_x2,
			cop.y2 as cover_partial_y2,
			qd.id as quad_dist_id,
			qd.macro_partial_id as quad_dist_macro_partial_id,
			qd.depth as quad_dist_depth,
			qd.area as quad_dist_area,
			qd.dist as quad_dist_dist
		from quad_dists qd
		inner join macro_partials map
			on qd.macro_partial_id = map.id
		inner join cover_partials cop
			on map.cover_partial_id = cop.id
		where map.macro_id = ?`)

	if depth > 0 {
		sqlStr = fmt.Sprintf("%s and qd.depth <= %d", sqlStr, depth)
	}

	if area > 0 {
		sqlStr = fmt.Sprintf("%s and qd.area >= %d", sqlStr, area)
	}

	sqlStr = fmt.Sprintf("%s order by qd.dist desc limit 1", sqlStr)

	var v model.CoverPartialQuadView = model.CoverPartialQuadView{
		CoverPartial: &model.CoverPartial{},
		QuadDist:     &model.QuadDist{},
	}

	err := s.dbMap.Db.QueryRow(sqlStr, macro.Id).Scan(
		&v.CoverPartial.Id,
		&v.CoverPartial.CoverId,
		&v.CoverPartial.AspectId,
		&v.CoverPartial.X1,
		&v.CoverPartial.Y1,
		&v.CoverPartial.X2,
		&v.CoverPartial.Y2,
		&v.QuadDist.Id,
		&v.QuadDist.MacroPartialId,
		&v.QuadDist.Depth,
		&v.QuadDist.Area,
		&v.QuadDist.Dist,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &v, nil
}
