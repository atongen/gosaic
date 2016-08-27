package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"log"
	"os"
	"os/signal"
	"sync"
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

			gidxs, err := gidxPartialService.FindMissing(aspect, "gidx.id ASC", 200, 0)
			if err != nil {
				return err
			}

			if len(gidxs) == 0 {
				break
			}

			gidxPartials := buildGidxPartials(l, gidxs, aspect)

			if len(gidxPartials) > 0 {
				numCreated, err := gidxPartialService.BulkInsert(gidxPartials)
				if err != nil {
					return err
				}
				bar.Add(int(numCreated))
			}
		}
	}

	bar.Finish()

	return nil
}

func buildGidxPartials(l *log.Logger, gidxs []*model.Gidx, aspect *model.Aspect) []*model.GidxPartial {
	var wg sync.WaitGroup
	wg.Add(len(gidxs))

	gidxPartials := make([]*model.GidxPartial, len(gidxs))
	add := make(chan *model.GidxPartial)
	errs := make(chan error)

	go func(myLog *log.Logger, gps []*model.GidxPartial, addCh chan *model.GidxPartial, errsCh chan error) {
		idx := 0
		for i := 0; i < len(gps); i++ {
			select {
			case gp := <-addCh:
				gps[idx] = gp
				idx++
			case err := <-errsCh:
				l.Printf("Error building index partial: %s\n", err.Error())
			}
			wg.Done()
		}
	}(l, gidxPartials, add, errs)

	for _, gidx := range gidxs {
		go func(g *model.Gidx, addCh chan *model.GidxPartial, errsCh chan error) {
			gp, err := buildGidxPartial(g, aspect)
			if err != nil {
				errsCh <- err
				return
			}
			addCh <- gp
		}(gidx, add, errs)
	}

	wg.Wait()
	close(add)
	close(errs)

	return gidxPartials
}

func buildGidxPartial(gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
	}

	pixels, err := util.GetAspectLab(gidx, aspect)
	if err != nil {
		return nil, err
	}
	p.Pixels = pixels

	err = p.EncodePixels()
	if err != nil {
		return nil, err
	}

	return &p, nil
}
