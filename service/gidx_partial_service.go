package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"sync"

	"gopkg.in/gorp.v1"
)

type GidxPartialService interface {
	Service
	Insert(*model.GidxPartial) error
	BulkInsert([]*model.GidxPartial) (int64, error)
	Update(*model.GidxPartial) error
	Delete(*model.GidxPartial) error
	Get(int64) (*model.GidxPartial, error)
	GetOneBy(string, interface{}) (*model.GidxPartial, error)
	ExistsBy(string, ...interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, ...interface{}) (int64, error)
	CountForMacro(*model.Macro) (int64, error)
	Find(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	Create(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindOrCreate(*model.Gidx, *model.Aspect) (*model.GidxPartial, error)
	FindMissing(*model.Aspect, string, int, int) ([]*model.Gidx, error)
	CountMissing([]*model.Aspect) (int64, error)
}

type gidxPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewGidxPartialService(dbMap *gorp.DbMap) GidxPartialService {
	return &gidxPartialServiceImpl{dbMap: dbMap}
}

func (s *gidxPartialServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.GidxPartial{}, "gidx_partials").SetKeys(true, "id")
	return nil
}

func (s *gidxPartialServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *gidxPartialServiceImpl) Insert(gidxPartial *model.GidxPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	err := gidxPartial.EncodePixels()
	if err != nil {
		return err
	}
	return s.dbMap.Insert(gidxPartial)
}

func (s *gidxPartialServiceImpl) BulkInsert(gidxPartials []*model.GidxPartial) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(gidxPartials) == 0 {
		return int64(0), nil
	} else if len(gidxPartials) == 1 {
		err := s.dbMap.Insert(gidxPartials[0])
		if err != nil {
			return int64(0), err
		}
		return int64(1), nil
	}

	var b bytes.Buffer

	b.WriteString("insert into gidx_partials (id, gidx_id, aspect_id, data) ")
	b.WriteString(fmt.Sprintf("select null as id, %d as gidx_id, %d as aspect_id, '%s' as data",
		gidxPartials[0].GidxId, gidxPartials[0].AspectId, gidxPartials[0].Data))

	for i := 1; i < len(gidxPartials); i++ {
		b.WriteString(fmt.Sprintf(" union select null, %d, %d, '%s'",
			gidxPartials[i].GidxId, gidxPartials[i].AspectId, gidxPartials[i].Data))
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

func (s *gidxPartialServiceImpl) Update(gidxPartial *model.GidxPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	err := gidxPartial.EncodePixels()
	if err != nil {
		return err
	}
	_, err = s.dbMap.Update(gidxPartial)
	return err
}

func (s *gidxPartialServiceImpl) Delete(gidxPartial *model.GidxPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(gidxPartial)
	return err
}

func (s *gidxPartialServiceImpl) Get(id int64) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	gidxPartial, err := s.dbMap.Get(model.GidxPartial{}, id)
	if err != nil {
		return nil, err
	} else if gidxPartial == nil {
		return nil, nil
	}

	gp, ok := gidxPartial.(*model.GidxPartial)
	if !ok {
		return nil, fmt.Errorf("Received struct is not a GidxPartial")
	}

	err = gp.DecodeData()
	if err != nil {
		return nil, err
	}

	return gp, nil
}

func (s *gidxPartialServiceImpl) GetOneBy(column string, value interface{}) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var gidxPartial model.GidxPartial
	err := s.dbMap.SelectOne(&gidxPartial, fmt.Sprintf("select * from gidx_partials where %s = ? limit 1", column), value)
	if err != nil {
		return nil, err
	}

	err = gidxPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	return &gidxPartial, err
}

func (s *gidxPartialServiceImpl) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	count, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from gidx_partials where %s limit 1", conditions), params...)
	return count == 1, err
}

func (s *gidxPartialServiceImpl) Count() (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from gidx_partials")
}

func (s *gidxPartialServiceImpl) CountBy(conditions string, params ...interface{}) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`
		select count(*)
		from gidx_partials
		where %s
		limit 1
	`, conditions)

	return s.dbMap.SelectInt(sql, params...)
}

func (s *gidxPartialServiceImpl) CountForMacro(macro *model.Macro) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := `
		select count(*)
		from gidx g
		where exists (
			select 1
			from gidx_partials gp,
			macro_partials mp
			where mp.macro_id = ?
			and mp.aspect_id = gp.aspect_id
			and gp.gidx_id = g.id
		);
	`
	return s.dbMap.SelectInt(sql, macro.Id)
}

func (s *gidxPartialServiceImpl) doFind(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	err := s.dbMap.SelectOne(&p, "select * from gidx_partials where gidx_id = ? and aspect_id = ? limit 1", p.GidxId, p.AspectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	err = p.DecodeData()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *gidxPartialServiceImpl) Find(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doFind(gidx, aspect)
}

func (s *gidxPartialServiceImpl) doCreate(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	pixels, err := util.GetAspectLab(gidx, aspect)
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

func (s *gidxPartialServiceImpl) Create(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.doCreate(gidx, aspect)
}

func (s *gidxPartialServiceImpl) FindOrCreate(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p, err := s.doFind(gidx, aspect)
	if err != nil {
		return nil, err
	} else if p != nil {
		return p, nil
	}

	// or create
	return s.doCreate(gidx, aspect)
}

func (s *gidxPartialServiceImpl) FindMissing(aspect *model.Aspect, order string, limit, offset int) ([]*model.Gidx, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf(`
		select *
		from gidx
		where not exists (
			select 1 from gidx_partials
			where gidx_partials.gidx_id = gidx.id
			and gidx_partials.aspect_id = ?
		)
		order by %s
		limit %d
		offset %d
	`, order, limit, offset)

	var gidxs []*model.Gidx
	_, err := s.dbMap.Select(&gidxs, sql, aspect.Id)

	return gidxs, err
}

func (s *gidxPartialServiceImpl) CountMissing(aspects []*model.Aspect) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	num := len(aspects)
	if num == 0 {
		return 0, nil
	}

	var count int64 = 0

	for _, aspect := range aspects {
		sql := `
			select count(*)
			from gidx
			where not exists (
				select 1 from gidx_partials
				where gidx_partials.gidx_id = gidx.id
				and gidx_partials.aspect_id = ?
			)
		`
		ac, err := s.dbMap.SelectInt(sql, aspect.Id)
		if err != nil {
			return int64(0), err
		}

		count += ac
	}

	return count, nil
}
