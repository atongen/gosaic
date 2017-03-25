package service

type Service interface {
	Register() error
	Close() error
}
