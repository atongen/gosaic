package controller

import (
	"errors"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"image"

	"gopkg.in/cheggaaa/pb.v1"
)

func PartialAspect(env environment.Environment, macroId int64, threashold float64) error {
	aspectService := env.ServiceFactory().MustAspectService()
	macroService := env.ServiceFactory().MustMacroService()
	macroPartialService := env.ServiceFactory().MustMacroPartialService()

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

	err = createPartialGidxIndexes(env, aspects, threashold, env.Workers())
	if err != nil {
		env.Printf("Error creating index aspects: %s\n", err.Error())
		return err
	}

	return nil
}

func createPartialGidxIndexes(env environment.Environment, aspects []*model.Aspect, threashold float64, workers int) error {
	gidxService := env.ServiceFactory().MustGidxService()
	gidxPartialService := env.ServiceFactory().MustGidxPartialService()

	count, err := gidxService.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	env.Printf("Building %d index image partials...\n", count)
	bar := pb.StartNew(int(count))
	batchSize := workers * 8

	for i := 0; ; i++ {
		gidxs, err := gidxService.FindAll("id asc", batchSize, i*batchSize)
		if err != nil {
			return err
		}

		if len(gidxs) == 0 {
			break
		}

		for _, gidx := range gidxs {
			if env.Cancel() {
				return errors.New("Cancelled")
			}

			gidxPartials, err := buildGidxPartials(env, gidx, aspects, threashold, workers)
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

func buildGidxPartials(env environment.Environment, gidx *model.Gidx, aspects []*model.Aspect, threashold float64, workers int) ([]*model.GidxPartial, error) {
	gidxPartialService := env.ServiceFactory().MustGidxPartialService()

	var gidxPartials []*model.GidxPartial

	pAspects := []*model.Aspect{}
	if threashold < 0.0 {
		pAspects = aspects
	} else {
		for _, aspect := range aspects {
			if gidx.Within(threashold, aspect) {
				pAspects = append(pAspects, aspect)
			}
		}
	}

	if len(pAspects) == 0 {
		return gidxPartials, nil
	}

	myAspects := []*model.Aspect{}
	for _, aspect := range pAspects {
		exists, err := gidxPartialService.ExistsBy("gidx_id = ? and aspect_id = ?", gidx.Id, aspect.Id)
		if err != nil {
			return nil, err
		} else if !exists {
			myAspects = append(myAspects, aspect)
		}

	}

	if len(myAspects) == 0 {
		return gidxPartials, nil
	}

	gidxPartials = make([]*model.GidxPartial, len(myAspects))

	img, err := util.OpenImg(gidx)
	if err != nil {
		return nil, err
	}

	add := make(chan *model.GidxPartial)
	errs := make(chan error)
	sem := make(chan bool, workers)

	go func(gps []*model.GidxPartial) {
		idx := 0
		for i := 0; i < len(myAspects); i++ {
			select {
			case gp := <-add:
				gps[idx] = gp
				idx += 1
			case err := <-errs:
				env.Printf("Error building index partial: %s\n", err.Error())
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
