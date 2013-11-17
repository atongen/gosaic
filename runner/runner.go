package runner

type Run struct {
	Path string
	Arg  string
}

type Runner interface {
	Execute() error
}
