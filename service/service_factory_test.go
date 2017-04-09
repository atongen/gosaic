package service

import (
	"testing"
)

func getServiceFactory() (ServiceFactory, error) {
	return NewServiceFactory("sqlite3://:memory:")
}

func TestServices(t *testing.T) {
	f, err := getServiceFactory()
	if err != nil {
		t.Fatalf("Error getting service factory: %s\n", err.Error())
	}
	defer f.Close()

	_, err = f.GidxService()
	if err != nil {
		t.Fatalf("Error getting gidxService: %s\n", err.Error())
	}

	_, err = f.AspectService()
	if err != nil {
		t.Fatalf("Error getting aspectService: %s\n", err.Error())
	}

	_, err = f.GidxPartialService()
	if err != nil {
		t.Fatalf("Error getting gidxPartialService: %s\n", err.Error())
	}

	_, err = f.CoverService()
	if err != nil {
		t.Fatalf("Error getting coverService: %s\n", err.Error())
	}

	_, err = f.CoverPartialService()
	if err != nil {
		t.Fatalf("Error getting coverService: %s\n", err.Error())
	}

	_, err = f.MacroService()
	if err != nil {
		t.Fatalf("Error getting macroService: %s\n", err.Error())
	}

	_, err = f.MacroPartialService()
	if err != nil {
		t.Fatalf("Error getting macroPartialService: %s\n", err.Error())
	}

	_, err = f.PartialComparisonService()
	if err != nil {
		t.Fatalf("Error getting partialComparisonService: %s\n", err.Error())
	}

	_, err = f.MosaicService()
	if err != nil {
		t.Fatalf("Error getting mosaicService: %s\n", err.Error())
	}

	_, err = f.MosaicPartialService()
	if err != nil {
		t.Fatalf("Error getting mosaicPartialService: %s\n", err.Error())
	}

	_, err = f.QuadDistService()
	if err != nil {
		t.Fatalf("Error getting quadDistService: %s\n", err.Error())
	}

	_, err = f.ProjectService()
	if err != nil {
		t.Fatalf("Error getting projectService: %s\n", err.Error())
	}
}

func TestMustServices(t *testing.T) {
	f, err := getServiceFactory()
	if err != nil {
		t.Fatalf("Error getting service factory: %s\n", err.Error())
	}
	defer f.Close()

	f.MustGidxService()
	f.MustAspectService()
	f.MustGidxPartialService()
	f.MustCoverService()
	f.MustCoverPartialService()
	f.MustMacroService()
	f.MustMacroPartialService()
	f.MustPartialComparisonService()
	f.MustMosaicService()
	f.MustMosaicPartialService()
	f.MustQuadDistService()
	f.MustProjectService()
}
