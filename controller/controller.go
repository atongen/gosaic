package controller

import (
	"github.com/atongen/gosaic/model"
)

type Executable struct {
	Project *model.Project
	Arg     string
}

type Controller interface {
	Execute() error
}
