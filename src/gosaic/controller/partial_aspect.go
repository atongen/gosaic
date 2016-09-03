package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func PartialAspect(env environment.Environment, macroId int64, threashold float64) error {
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

	gidxService, err := env.GidxService()
	if err != nil {
		env.Printf("Error getting gidx service: %s\n", err.Error())
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

	err = createPartialGidxIndexes(env.Log(), gidxService, gidxPartialService, aspects, threashold, env.Workers())
	if err != nil {
		env.Printf("Error creating index aspects: %s\n", err.Error())
		return err
	}

	return nil
}

func createPartialGidxIndexes(l *log.Logger, gidxService service.GidxService, gidxPartialService service.GidxPartialService, aspects []*model.Aspect, threashold float64, workers int) error {
	count, err := gidxService.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	l.Printf("Building %d index image partials...\n", count)
	bar := pb.StartNew(int(count))
	batchSize := workers * 8

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for i := 0; ; i++ {
		gidxs, err := gidxService.FindAll("id asc", batchSize, i*batchSize)
		if err != nil {
			return err
		}

		if len(gidxs) == 0 {
			break
		}

		for _, gidx := range gidxs {
			if cancel {
				return errors.New("Cancelled")
			}

			gidxPartials, err := buildGidxPartials(l, gidxPartialService, gidx, aspects, threashold, workers)
			if err != nil {
				return err
			}

			_, err = gidxPartialService.BulkInsert(gidxPartials)
			if err != nil {
				return err
			}
			bar.Increment()
		}
	}

	bar.Finish()

	return nil
}

func buildGidxPartials(l *log.Logger, gidxPartialService service.GidxPartialService, gidx *model.Gidx, aspects []*model.Aspect, threashold float64, workers int) ([]*model.GidxPartial, error) {
	myAspects := []*model.Aspect{}
	if threashold < 0.0 {
		myAspects = aspects
	} else {
		for _, aspect := range aspects {
			if gidx.Within(threashold, aspect) {
				exists, err := gidxPartialService.ExistsBy("gidx_id = ? and aspect_id = ?", gidx.Id, aspect.Id)
				if err != nil {
					return nil, err
				} else if !exists {
					myAspects = append(myAspects, aspect)
				}
			}
		}
	}
	num := len(myAspects)

	var gidxPartials []*model.GidxPartial
	if num == 0 {
		return gidxPartials, nil
	}

	gidxPartials = make([]*model.GidxPartial, num)

	img, err := util.OpenImg(gidx)
	if err != nil {
		return nil, err
	}

	add := make(chan *model.GidxPartial)
	errs := make(chan error)
	sem := make(chan bool, workers)

	go func(gps []*model.GidxPartial) {
		idx := 0
		for i := 0; i < num; i++ {
			select {
			case gp := <-add:
				gps[idx] = gp
				idx += 1
			case err := <-errs:
				l.Printf("Error building index partial: %s\n", err.Error())
			}
			<-sem
		}
	}(gidxPartials)

	for _, aspect := range myAspects {
		sem <- true
		go func(a *model.Aspect) {
			gp, err := buildGidxPartial(img, gidx, a)
			if err != nil {
				errs <- err
				return
			}
			add <- gp
		}(aspect)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(add)
	close(errs)
	close(sem)

	return gidxPartials, nil
}

func buildGidxPartial(img *image.Image, gidx *model.Gidx, aspect *model.Aspect) (*model.GidxPartial, error) {
	p := model.GidxPartial{
		GidxId:   gidx.Id,
		AspectId: aspect.Id,
		Pixels:   util.GetImgAspectLab(img, gidx, aspect),
	}

	err := p.EncodePixels()
	if err != nil {
		return nil, err
	}

	return &p, nil
}
