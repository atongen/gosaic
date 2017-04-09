package service

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/atongen/gosaic/database"
	gorp "gopkg.in/gorp.v1"

	_ "github.com/mattn/go-sqlite3"
)

type ServiceName uint8

const (
	GidxServiceName ServiceName = iota
	AspectServiceName
	GidxPartialServiceName
	CoverServiceName
	CoverPartialServiceName
	MacroServiceName
	MacroPartialServiceName
	PartialComparisonServiceName
	MosaicServiceName
	MosaicPartialServiceName
	QuadDistServiceName
	ProjectServiceName
)

type ServiceFactory interface {
	Close() error

	GidxService() (GidxService, error)
	AspectService() (AspectService, error)
	GidxPartialService() (GidxPartialService, error)
	CoverService() (CoverService, error)
	CoverPartialService() (CoverPartialService, error)
	MacroService() (MacroService, error)
	MacroPartialService() (MacroPartialService, error)
	PartialComparisonService() (PartialComparisonService, error)
	MosaicService() (MosaicService, error)
	MosaicPartialService() (MosaicPartialService, error)
	QuadDistService() (QuadDistService, error)
	ProjectService() (ProjectService, error)

	MustGidxService() GidxService
	MustAspectService() AspectService
	MustGidxPartialService() GidxPartialService
	MustCoverService() CoverService
	MustCoverPartialService() CoverPartialService
	MustMacroService() MacroService
	MustMacroPartialService() MacroPartialService
	MustPartialComparisonService() PartialComparisonService
	MustMosaicService() MosaicService
	MustMosaicPartialService() MosaicPartialService
	MustQuadDistService() QuadDistService
	MustProjectService() ProjectService
}

func NewServiceFactory(dsn string) (ServiceFactory, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	default:
		return nil, fmt.Errorf("Unknown database type: %s", u.Scheme)
	case "sqlite3":
		return newServiceFactorySqlite3(u)
	}
}

func newServiceFactorySqlite3(u *url.URL) (*serviceFactorySqlite3, error) {
	f := serviceFactorySqlite3{
		services: make(map[ServiceName]Service),
	}

	db, err := sql.Open("sqlite3", u.Path)
	if err != nil {
		return nil, err
	}
	f.dB = db

	err = f.dB.Ping()
	if err != nil {
		return nil, err
	}

	_, err = f.dB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	_, err = database.Migrate(f.dB)
	if err != nil {
		return nil, err
	}

	// setup orm
	f.dbMap = &gorp.DbMap{Db: f.dB, Dialect: gorp.SqliteDialect{}}

	return &f, nil
}
