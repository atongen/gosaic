package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/atongen/gosaic/model"

	"gopkg.in/gorp.v1"
)

type projectServiceSqlite3 struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewProjectService(dbMap *gorp.DbMap) *projectServiceSqlite3 {
	return &projectServiceSqlite3{dbMap: dbMap}
}

func (s *projectServiceSqlite3) Register() error {
	s.dbMap.AddTableWithName(model.Project{}, "projects").SetKeys(true, "id")
	return nil
}

func (s *projectServiceSqlite3) Close() error {
	return s.dbMap.Db.Close()
}

func (s *projectServiceSqlite3) Get(id int64) (*model.Project, error) {
	s.m.Lock()
	defer s.m.Unlock()

	c, err := s.dbMap.Get(model.Project{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	m, ok := c.(*model.Project)
	if !ok {
		return nil, errors.New("Unable to type cast project")
	}

	if m.Id == int64(0) {
		return nil, nil
	}

	return m, nil
}

func (s *projectServiceSqlite3) Insert(project *model.Project) error {
	s.m.Lock()
	defer s.m.Unlock()

	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	return s.dbMap.Insert(project)
}

func (s *projectServiceSqlite3) Update(project *model.Project) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(project)
}

func (s *projectServiceSqlite3) GetOneBy(conditions string, params ...interface{}) (*model.Project, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var project model.Project

	err := s.dbMap.SelectOne(&project, fmt.Sprintf("select * from projects where %s limit 1", conditions), params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &project, nil
}

func (s *projectServiceSqlite3) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	s.m.Lock()
	defer s.m.Unlock()

	num, err := s.dbMap.SelectInt(fmt.Sprintf("select 1 from projects where %s limit 1", conditions), params...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}

	return num == 1, nil
}

func (s *projectServiceSqlite3) FindAll(order string) ([]*model.Project, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from projects order by %s", order)

	var projects []*model.Project
	_, err := s.dbMap.Select(&projects, sql, order)

	return projects, err
}
