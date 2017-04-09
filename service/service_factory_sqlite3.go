package service

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/atongen/gosaic/service/sqlite3"

	gorp "gopkg.in/gorp.v1"
)

type serviceFactorySqlite3 struct {
	m        sync.Mutex
	dB       *sql.DB
	dbMap    *gorp.DbMap
	services map[ServiceName]Service
}

func (f *serviceFactorySqlite3) getService(name ServiceName) (Service, error) {
	f.m.Lock()
	defer f.m.Unlock()

	var s Service
	if s, ok := f.services[name]; ok {
		return s, nil
	}

	switch name {
	default:
		return nil, fmt.Errorf("Service not found")
	case GidxServiceName:
		s = sqlite3.NewGidxService(f.dbMap)
	case AspectServiceName:
		s = sqlite3.NewAspectService(f.dbMap)
	case GidxPartialServiceName:
		s = sqlite3.NewGidxPartialService(f.dbMap)
	case CoverServiceName:
		s = sqlite3.NewCoverService(f.dbMap)
	case CoverPartialServiceName:
		s = sqlite3.NewCoverPartialService(f.dbMap)
	case MacroServiceName:
		s = sqlite3.NewMacroService(f.dbMap)
	case MacroPartialServiceName:
		s = sqlite3.NewMacroPartialService(f.dbMap)
	case PartialComparisonServiceName:
		s = sqlite3.NewPartialComparisonService(f.dbMap)
	case MosaicServiceName:
		s = sqlite3.NewMosaicService(f.dbMap)
	case MosaicPartialServiceName:
		s = sqlite3.NewMosaicPartialService(f.dbMap)
	case QuadDistServiceName:
		s = sqlite3.NewQuadDistService(f.dbMap)
	case ProjectServiceName:
		s = sqlite3.NewProjectService(f.dbMap)
	}

	err := s.Register()
	if err != nil {
		return nil, err
	}

	f.services[name] = s
	return s, nil
}

func (f *serviceFactorySqlite3) Close() error {
	return f.dB.Close()
}

func (f *serviceFactorySqlite3) GidxService() (GidxService, error) {
	s, err := f.getService(GidxServiceName)
	if err != nil {
		return nil, err
	}

	gidxService, ok := s.(GidxService)
	if !ok {
		return nil, fmt.Errorf("Invalid gidx service")
	}

	return gidxService, nil
}

func (f *serviceFactorySqlite3) AspectService() (AspectService, error) {
	s, err := f.getService(AspectServiceName)
	if err != nil {
		return nil, err
	}

	aspectService, ok := s.(AspectService)
	if !ok {
		return nil, fmt.Errorf("Invalid aspect service")
	}

	return aspectService, nil
}

func (f *serviceFactorySqlite3) GidxPartialService() (GidxPartialService, error) {
	s, err := f.getService(GidxPartialServiceName)
	if err != nil {
		return nil, err
	}

	gidxPartialService, ok := s.(GidxPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid gidx_partial service")
	}

	return gidxPartialService, nil
}

func (f *serviceFactorySqlite3) CoverService() (CoverService, error) {
	s, err := f.getService(CoverServiceName)
	if err != nil {
		return nil, err
	}

	coverService, ok := s.(CoverService)
	if !ok {
		return nil, fmt.Errorf("Invalid cover service")
	}

	return coverService, nil
}

func (f *serviceFactorySqlite3) CoverPartialService() (CoverPartialService, error) {
	s, err := f.getService(CoverPartialServiceName)
	if err != nil {
		return nil, err
	}

	coverPartialService, ok := s.(CoverPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid cover_partial service")
	}

	return coverPartialService, nil
}

func (f *serviceFactorySqlite3) MacroService() (MacroService, error) {
	s, err := f.getService(MacroServiceName)
	if err != nil {
		return nil, err
	}

	macroService, ok := s.(MacroService)
	if !ok {
		return nil, fmt.Errorf("Invalid macro service")
	}

	return macroService, nil
}

func (f *serviceFactorySqlite3) MacroPartialService() (MacroPartialService, error) {
	s, err := f.getService(MacroPartialServiceName)
	if err != nil {
		return nil, err
	}

	macroPartialService, ok := s.(MacroPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid macro partial service")
	}

	return macroPartialService, nil
}

func (f *serviceFactorySqlite3) PartialComparisonService() (PartialComparisonService, error) {
	s, err := f.getService(PartialComparisonServiceName)
	if err != nil {
		return nil, err
	}

	partialComparisonService, ok := s.(PartialComparisonService)
	if !ok {
		return nil, fmt.Errorf("Invalid partial comparison service")
	}

	return partialComparisonService, nil
}

func (f *serviceFactorySqlite3) MosaicService() (MosaicService, error) {
	s, err := f.getService(MosaicServiceName)
	if err != nil {
		return nil, err
	}

	mosaicService, ok := s.(MosaicService)
	if !ok {
		return nil, fmt.Errorf("Invalid mosaic service")
	}

	return mosaicService, nil
}

func (f *serviceFactorySqlite3) MosaicPartialService() (MosaicPartialService, error) {
	s, err := f.getService(MosaicPartialServiceName)
	if err != nil {
		return nil, err
	}

	mosaicPartialService, ok := s.(MosaicPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid mosaic partial service")
	}

	return mosaicPartialService, nil
}

func (f *serviceFactorySqlite3) QuadDistService() (QuadDistService, error) {
	s, err := f.getService(QuadDistServiceName)
	if err != nil {
		return nil, err
	}

	quadDistService, ok := s.(QuadDistService)
	if !ok {
		return nil, fmt.Errorf("Invalid quad dist service")
	}

	return quadDistService, nil
}

func (f *serviceFactorySqlite3) ProjectService() (ProjectService, error) {
	s, err := f.getService(ProjectServiceName)
	if err != nil {
		return nil, err
	}

	projectService, ok := s.(ProjectService)
	if !ok {
		return nil, fmt.Errorf("Invalid project service")
	}

	return projectService, nil
}

func (f *serviceFactorySqlite3) MustGidxService() GidxService {
	s, err := f.GidxService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustAspectService() AspectService {
	s, err := f.AspectService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustGidxPartialService() GidxPartialService {
	s, err := f.GidxPartialService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustCoverService() CoverService {
	s, err := f.CoverService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustCoverPartialService() CoverPartialService {
	s, err := f.CoverPartialService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustMacroService() MacroService {
	s, err := f.MacroService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustMacroPartialService() MacroPartialService {
	s, err := f.MacroPartialService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustPartialComparisonService() PartialComparisonService {
	s, err := f.PartialComparisonService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustMosaicService() MosaicService {
	s, err := f.MosaicService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustMosaicPartialService() MosaicPartialService {
	s, err := f.MosaicPartialService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustQuadDistService() QuadDistService {
	s, err := f.QuadDistService()
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (f *serviceFactorySqlite3) MustProjectService() ProjectService {
	s, err := f.ProjectService()
	if err != nil {
		panic(err.Error())
	}
	return s
}
