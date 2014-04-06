package service

import "github.com/coopernurse/gorp"

type Service interface {
	DbMap() *gorp.DbMap
	Register()
	TableName() string
	PrimaryKey() string
}
