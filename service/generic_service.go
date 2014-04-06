package service

import "github.com/coopernurse/gorp"

type GenericService interface {
	Insert(...interface{}) error
	Update(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Get(int64) (interface{}, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
}

type GenericServiceImpl struct {
	dbMap      *gorp.DbMap
	holder     interface{}
	tableName  string
	primaryKey string
}

func NewGenericService(dbMap *gorp.DbMap, holder interface{}, tableName string, primaryKey string) *GenericServiceImpl {
	return &GenericServiceImpl{
		dbMap:      dbMap,
		holder:     holder,
		tableName:  tableName,
		primaryKey: primaryKey}
}

func (s *GenericServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

func (s *GenericServiceImpl) Register() {
	s.DbMap().AddTableWithName(s.holder, s.tableName).SetKeys(true, s.primaryKey)
}

func (s *GenericServiceImpl) Insert(generics ...interface{}) error {
	return s.DbMap().Insert(generics...)
}

func (s *GenericServiceImpl) Update(generics ...interface{}) (int64, error) {
	return s.DbMap().Update(generics...)
}

func (s *GenericServiceImpl) Delete(generics ...interface{}) (int64, error) {
	return s.DbMap().Delete(generics...)
}

func (s *GenericServiceImpl) Get(id int64) (interface{}, error) {
	return s.DbMap().Get(s.holder, id)
}

func (s *GenericServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from "+s.tableName+" where "+column+" = ?", value)
	return count == 1, err
}

func (s *GenericServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from " + s.tableName)
}

func (s *GenericServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from "+s.tableName+" where "+column+" = ?", value)
}
