package environment

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atongen/gosaic/service"
)

var (
	Version   = "none"
	BuildTime = "none"
	BuildUser = "none"
	BuildHash = "none"
)

type Environment interface {
	Init() error
	Close()
	Cancel() bool
	Workers() int
	Log() *log.Logger
	ServiceFactory() service.ServiceFactory
	ProjectId() int64
	SetProjectId(id int64)
	Printf(format string, a ...interface{})
	Println(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalln(a ...interface{})
}

type environment struct {
	workers        int
	projectId      int64
	log            *log.Logger
	cancel         bool
	cancelCh       chan os.Signal
	serviceFactory service.ServiceFactory
}

func NewEnvironment(dsn string, out io.Writer, workers int) (Environment, error) {
	env := &environment{}

	// setup the environment logger
	env.log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime)
	env.cancel = false

	serviceFactory, err := service.NewServiceFactory(dsn)
	if err != nil {
		return nil, err
	}

	env.serviceFactory = serviceFactory
	env.workers = workers

	return env, nil
}

func GetProdEnv(dsn string, workers int) (Environment, error) {
	return NewEnvironment(dsn, os.Stdout, workers)
}

func GetTestEnv(out io.Writer) (Environment, error) {
	return NewEnvironment("sqlite3://:memory:", out, 2)
}

func (env *environment) Init() error {
	// cancel
	env.cancelCh = make(chan os.Signal, 2)
	signal.Notify(env.cancelCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-env.cancelCh
		env.Println("Interrupt caught, attempting graceful shutdown...")
		env.cancel = true
	}()

	return nil
}

func (env *environment) Close() {
	env.serviceFactory.Close()
	close(env.cancelCh)
}

func (env *environment) Cancel() bool {
	return env.cancel
}

func (env *environment) Workers() int {
	return env.workers
}

func (env *environment) Log() *log.Logger {
	return env.log
}

func (env *environment) ServiceFactory() service.ServiceFactory {
	return env.serviceFactory
}

func (env *environment) ProjectId() int64 {
	return env.projectId
}

func (env *environment) SetProjectId(projectId int64) {
	env.projectId = projectId
}

func (env *environment) Fatalln(v ...interface{}) {
	env.log.Fatalln(v...)
}

func (env *environment) Fatalf(format string, v ...interface{}) {
	env.log.Fatalf(format, v...)
}

func (env *environment) Println(v ...interface{}) {
	env.log.Println(v...)
}

func (env *environment) Printf(format string, v ...interface{}) {
	env.log.Printf(format, v...)
}
