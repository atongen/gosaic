package runner

import (
	"github.com/atongen/gosaic"
)

type Run struct {
	Project *gosaic.Project
	Arg     string
}

type Runner interface {
	Execute() error
}
