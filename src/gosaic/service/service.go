package service

import "gopkg.in/gorp.v1"

type Service interface {
	DbMap() *gorp.DbMap
	Register() error
}
