package service

import (
	"errors"
	"fmt"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

type MacroService interface {
	Service
	Get(int64) (*model.Macro, error)
	Insert(*model.Macro) error
	Update(*model.Macro) error
	Delete(*model.Macro) error
	GetOneBy(string, ...interface{}) (*model.Macro, error)
	ExistsBy(string, ...interface{}) (bool, error)
	FindAll(string) ([]*model.Macro, error)
}

type macroServiceImpl struct {
	dbMap *gorp.DbMap
	m     sync.Mutex
}

func NewMacroService(dbMap *gorp.DbMap) MacroService {
	return &macroServiceImpl{dbMap: dbMap}
}

func (s *macroServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *macroServiceImpl) Register() error {
	s.DbMap().AddTableWithName(model.Macro{}, "macros").SetKeys(true, "id")
	return nil
}

func (s *macroServiceImpl) Get(id int64) (*model.Macro, error) {
	c, err := s.DbMap().Get(model.Macro{}, id)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	m, ok := c.(*model.Macro)
	if !ok {
		return nil, errors.New("Unable to type cast macro")
	}

	if m.Id == int64(0) {
		return nil, nil
	}

	return m, nil
}

func (s *macroServiceImpl) Insert(c *model.Macro) error {
	return s.DbMap().Insert(c)
}

func (s *macroServiceImpl) Update(c *model.Macro) error {
	_, err := s.DbMap().Update(c)
	return err
}

func (s *macroServiceImpl) Delete(c *model.Macro) error {
	_, err := s.DbMap().Delete(c)
	return err
}

func (s *macroServiceImpl) GetOneBy(conditions string, params ...interface{}) (*model.Macro, error) {
	var macro model.Macro

	err := s.DbMap().SelectOne(&macro, fmt.Sprintf("select * from macros where %s limit 1", conditions), params...)
	// returns error if none are found
	// or if more than one is found
	if err != nil {
		return nil, nil
	}

	return &macro, nil
}

func (s *macroServiceImpl) ExistsBy(conditions string, params ...interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt(fmt.Sprintf("select 1 from macros where %s limit 1", conditions), params...)
	return count == 1, err
}

func (s *macroServiceImpl) FindAll(order string) ([]*model.Macro, error) {
	sql := `select * from macros order by ?`

	var macros []*model.Macro
	_, err := s.dbMap.Select(&macros, sql, order)

	return macros, err
}
