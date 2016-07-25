package service

import (
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
	GetOneBy(string, interface{}) (*model.Macro, error)
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
	} else if c != nil {
		return c.(*model.Macro), nil
	} else {
		return nil, nil
	}
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

func (s *macroServiceImpl) GetOneBy(column string, value interface{}) (*model.Macro, error) {
	var macro model.Macro
	err := s.DbMap().SelectOne(&macro, "select * from macros where "+column+" = ? limit 1", value)
	return &macro, err
}

func (s *macroServiceImpl) FindAll(order string) ([]*model.Macro, error) {
	sql := `select * from macros order by ?`

	var macros []*model.Macro
	_, err := s.dbMap.Select(&macros, sql, order)

	return macros, err
}
