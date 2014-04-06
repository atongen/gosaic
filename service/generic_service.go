package service

import (
	"github.com/coopernurse/gorp"

	"github.com/atongen/gosaic/model"
)

type GenericService interface {
	Insert(...interface{}) error
	Update(...interface{}) (int64, error)
	Delete(...interface{}) (int64, error)
	Get(int64) (interface{}, error)
	GetOneBy(string, interface{}) (interface{}, error)
	ExistsBy(string, interface{}) (bool, error)
	Count() (int64, error)
	CountBy(string, interface{}) (int64, error)
}

type GenericServiceImpl struct {
	dbMap      *gorp.DbMap
	tableName  string
	primaryKey string
}

func NewGenericService(dbMap *gorp.DbMap, tableName string, primaryKey string) *GenericServiceImpl {
	return &GenericServiceImpl{dbMap: dbMap, tableName: tableName, primaryKey: primaryKey}
}

func (s *GenericServiceImpl) DbMap() *gorp.DbMap {
	return s.dbMap
}

// TODO: Reflection
func (s *GenericServiceImpl) Register() {
	s.DbMap().AddTableWithName(model.Gidx{}, s.TableName()).SetKeys(true, s.PrimaryKey())
}

func (s *GenericServiceImpl) TableName() string {
	return s.tableName
}

func (s *GenericServiceImpl) PrimaryKey() string {
	return s.primaryKey
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

// TODO: Reflection
func (s *GenericServiceImpl) Get(id int64) (interface{}, error) {
	return s.DbMap().Get(model.Gidx{}, id)
}

func (s *GenericServiceImpl) GetOneBy(column string, value interface{}) (interface{}, error) {
	var generic interface{}
	err := s.DbMap().SelectOne(&generic, "select * from "+s.TableName()+" where "+column+" = ?", value)
	return generic, err
}

func (s *GenericServiceImpl) ExistsBy(column string, value interface{}) (bool, error) {
	count, err := s.DbMap().SelectInt("select 1 from "+s.TableName()+" where "+column+" = ?", value)
	return count == 1, err
}

func (s *GenericServiceImpl) Count() (int64, error) {
	return s.DbMap().SelectInt("select count(*) from " + s.TableName())
}

func (s *GenericServiceImpl) CountBy(column string, value interface{}) (int64, error) {
	return s.DbMap().SelectInt("select count(*) from "+s.TableName()+" where "+column+" = ?", value)
}
