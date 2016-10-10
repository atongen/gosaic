package service

import (
	"database/sql"
	"errors"
	"fmt"
	"gosaic/model"
	"sync"
	"time"

	"gopkg.in/gorp.v1"
)

type ProjectService interface {
	Service
	Get(int64) (*model.Project, error)
	Insert(*model.Project) error
	Update(*model.Project) (int64, error)
	GetOneBy(string, ...interface{}) (*model.Project, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Project, error)
}

type projectServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewProjectService(dbMap *gorp.DbMap) ProjectService {
	return &projectServiceImpl{dbMap: dbMap}
}

func (s *projectServiceImpl) Register() error {
	s.dbMap.AddTableWithName(model.Project{}, "projects").SetKeys(true, "id")
	return nil
}

func (s *projectServiceImpl) Close() error {
	return s.dbMap.Db.Close()
}

func (s *projectServiceImpl) Get(id int64) (*model.Project, error) {
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

func (s *projectServiceImpl) Insert(project *model.Project) error {
	s.m.Lock()
	defer s.m.Unlock()

	if project.CreatedAt.IsZero() {
		project.CreatedAt = time.Now()
	}
	return s.dbMap.Insert(project)
}

func (s *projectServiceImpl) Update(project *model.Project) (int64, error) {
	s.m.Lock()
	defer s.m.Unlock()

	return s.dbMap.Update(project)
}

func (s *projectServiceImpl) GetOneBy(conditions string, params ...interface{}) (*model.Project, error) {
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

func (s *projectServiceImpl) ExistsBy(conditions string, params ...interface{}) (bool, error) {
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

func (s *projectServiceImpl) FindAll(order string) ([]*model.Project, error) {
	s.m.Lock()
	defer s.m.Unlock()

	sql := fmt.Sprintf("select * from projects order by %s", order)

	var projects []*model.Project
	_, err := s.dbMap.Select(&projects, sql, order)

	return projects, err
}
