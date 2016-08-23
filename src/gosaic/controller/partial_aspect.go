package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func PartialAspect(env environment.Environment, macroId int64) error {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return err
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
		return err
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error getting macro partial service: %s\n", err.Error())
		return err
	}

	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		env.Printf("Error getting gidx partial service: %s\n", err.Error())
		return err
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Printf("Error getting macro %d: %s\n", macroId, err.Error())
		return err
	}

	aspectIds, err := macroPartialService.AspectIds(macro.Id)
	if err != nil {
		env.Printf("Error getting aspect ids: %s\n", err.Error())
		return err
	}

	if len(aspectIds) == 0 {
		return nil
	}

	aspects, err := aspectService.FindIn(aspectIds)
	if err != nil {
		env.Printf("Error getting aspects: %s\n", err.Error())
		return err
	}

	err = createMissingGidxIndexes(env.Log(), gidxPartialService, aspects)
	if err != nil {
		env.Printf("Error creating index aspects: %s\n", err.Error())
		return err
	}

	return nil
}

func createMissingGidxIndexes(l *log.Logger, gidxPartialService service.GidxPartialService, aspects []*model.Aspect) error {
	count, err := gidxPartialService.CountMissing(aspects)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	l.Printf("Building %d indexed image partials...\n", count)
	bar := pb.StartNew(int(count))

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for _, aspect := range aspects {
		for {
			if cancel {
				return errors.New("Cancelled")
			}

			gidxs, err := gidxPartialService.FindMissing(aspect, "gidx.id ASC", 100, 0)
			if err != nil {
				return err
			}

			if len(gidxs) == 0 {
				break
			}

			for _, gidx := range gidxs {
				_, err := gidxPartialService.Create(gidx, aspect)
				if err != nil {
					return err
				}
				bar.Increment()
			}
		}
	}

	bar.Finish()

	return nil
}
