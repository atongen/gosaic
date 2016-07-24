package environment

import (
	"fmt"
	"gosaic/service"
)

type ServiceName uint8

const (
	GidxServiceName ServiceName = iota
	AspectServiceName
	GidxPartialServiceName
	CoverServiceName
	CoverPartialServiceName
)

func (env *environment) getService(name ServiceName) (service.Service, error) {
	env.m.Lock()
	defer env.m.Unlock()

	var s service.Service
	if s, ok := env.services[name]; ok {
		return s, nil
	}

	switch name {
	default:
		return nil, fmt.Errorf("Service not found")
	case GidxServiceName:
		s = service.NewGidxService(env.dbMap)
	case AspectServiceName:
		s = service.NewAspectService(env.dbMap)
	case GidxPartialServiceName:
		s = service.NewGidxPartialService(env.dbMap)
	case CoverServiceName:
		s = service.NewCoverService(env.dbMap)
	case CoverPartialServiceName:
		s = service.NewCoverPartialService(env.dbMap)
	}
	err := s.Register()
	if err != nil {
		return nil, err
	}
	env.services[name] = s
	return s, nil
}

func (env *environment) GidxService() (service.GidxService, error) {
	s, err := env.getService(GidxServiceName)
	if err != nil {
		return nil, err
	}

	gidxService, ok := s.(service.GidxService)
	if !ok {
		return nil, fmt.Errorf("Invalid gidx service")
	}

	return gidxService, nil
}

func (env *environment) AspectService() (service.AspectService, error) {
	s, err := env.getService(AspectServiceName)
	if err != nil {
		return nil, err
	}

	aspectService, ok := s.(service.AspectService)
	if !ok {
		return nil, fmt.Errorf("Invalid aspect service")
	}

	return aspectService, nil
}

func (env *environment) GidxPartialService() (service.GidxPartialService, error) {
	s, err := env.getService(GidxPartialServiceName)
	if err != nil {
		return nil, err
	}

	gidxPartialService, ok := s.(service.GidxPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid gidx_partial service")
	}

	return gidxPartialService, nil
}

func (env *environment) CoverService() (service.CoverService, error) {
	s, err := env.getService(CoverServiceName)
	if err != nil {
		return nil, err
	}

	coverService, ok := s.(service.CoverService)
	if !ok {
		return nil, fmt.Errorf("Invalid cover service")
	}

	return coverService, nil
}

func (env *environment) CoverPartialService() (service.CoverPartialService, error) {
	s, err := env.getService(CoverPartialServiceName)
	if err != nil {
		return nil, err
	}

	coverPartialService, ok := s.(service.CoverPartialService)
	if !ok {
		return nil, fmt.Errorf("Invalid cover_partial service")
	}

	return coverPartialService, nil
}
