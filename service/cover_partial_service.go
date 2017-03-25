package service

import (
	"bytes"
	"fmt"
	"github.com/atongen/gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type CoverPartialService interface {
	Service
	Get(int64) (*model.CoverPartial, error)
	Insert(*model.CoverPartial) error
	BulkInsert([]*model.CoverPartial) (int64, error)
	Count(*model.Cover) (int64, error)
	Update(*model.CoverPartial) error
	Delete(*model.CoverPartial) error
	FindAll(int64, string) ([]*model.CoverPartial, error)
}

type coverPartialServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewCoverPartialService(dbMap *gorp.DbMap) CoverPartialService {
	return &coverPartialServiceImpl{dbMap: dbMap}
}

func (s *coverPartialServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.CoverPartial{}, "cover_partials").SetKeys(true, "id")
	return nil
}

func (s *coverPartialServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *coverPartialServiceImpl) Get(id int64) (*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.CoverPartial{}, id)
	if err != nil {
		return nil, err
	} else if c != nil {
		return c.(*model.CoverPartial), nil
	} else {
		return nil, nil
	}
}

func (s *coverPartialServiceImpl) Insert(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Insert(c)
}

func (s *coverPartialServiceImpl) BulkInsert(coverPartials []*model.CoverPartial) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if len(coverPartials) == 0 {
		return int64(0), nil
	} else if len(coverPartials) == 1 {
		err := s.dbMap.Insert(coverPartials[0])
		if err != nil {
			return int64(0), err
		}
		return int64(1), nil
	}

	var b bytes.Buffer

	b.WriteString("insert into cover_partials (id, cover_id, aspect_id, x1, y1, x2, y2) ")
	b.WriteString(fmt.Sprintf("select null as id, %d as cover_id, %d as aspect_id, %d as x1, %d as y1, %d as x2, %d as y2",
		coverPartials[0].CoverId, coverPartials[0].AspectId, coverPartials[0].X1, coverPartials[0].Y1, coverPartials[0].X2, coverPartials[0].Y2))

	for i := 1; i < len(coverPartials); i++ {
		b.WriteString(fmt.Sprintf(" union select null, %d, %d, %d, %d, %d, %d",
			coverPartials[i].CoverId, coverPartials[i].AspectId, coverPartials[i].X1, coverPartials[i].Y1, coverPartials[i].X2, coverPartials[i].Y2))
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

func (s *coverPartialServiceImpl) Count(c *model.Cover) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.SelectInt("select count(*) from cover_partials where cover_id = ? limit 1", c.Id)
}

func (s *coverPartialServiceImpl) Update(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Update(c)
	return err
}

func (s *coverPartialServiceImpl) Delete(c *model.CoverPartial) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.dbMap.Delete(c)
	return err
}

func (s *coverPartialServiceImpl) FindAll(coverId int64, order string) ([]*model.CoverPartial, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from cover_partials where cover_id = ? order by %s", order)

	var coverPartials []*model.CoverPartial
	_, err := s.dbMap.Select(&coverPartials, sql, coverId)

	return coverPartials, err
}
