package controller

import (
	"errors"
	"fmt"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"image"
	"math"
	"strings"

	"gopkg.in/cheggaaa/pb.v1"
)

func MacroQuad(env environment.Environment,
	path string,
	coverWidth, coverHeight, size, minDepth, maxDepth, minArea, maxArea int,
	coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {

	aspectService := env.ServiceFactory().MustAspectService()
	coverService := env.ServiceFactory().MustCoverService()

	myCoverWidth, myCoverHeight, err := calculateDimensions(aspectService, path, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover dimensions: %s\n", err.Error())
		return nil, nil
	}

	size, minDepth, maxDepth, minArea, maxArea, err = macroQuadFixArgs(myCoverWidth, myCoverHeight, size, minDepth, maxDepth, minArea, maxArea)
	if err != nil {
		env.Printf("Error calculating quad arguments: %s\n", err.Error())
		return nil, nil
	}

	cover, err := envCover(env)
	if err != nil {
		env.Printf("Error getting cover from project environment: %s\n", err.Error())
		return nil, nil
	}

	if cover == nil {
		cover, err = macroQuadCreateCover(env, myCoverWidth, myCoverHeight)
		if err != nil {
			env.Printf("Error building cover: %s\n", err.Error())
			return nil, nil
		}
	}

	err = setEnvCover(env, cover)
	if err != nil {
		env.Printf("Error setting cover in project environment: %s\n", err.Error())
		return nil, nil
	}

	macro, img, err := findOrCreateMacro(env, cover, path, macroOutfile)
	if err != nil {
		env.Printf("Error building macro: %s\n", err.Error())
		coverService.Delete(cover)
		return cover, nil
	}

	err = setEnvMacro(env, macro)
	if err != nil {
		env.Printf("Error setting macro in project environment: %s\n", err.Error())
		return cover, nil
	}

	err = macroQuadBuildPartials(env, cover, macro, img, size, minDepth, maxDepth, minArea, maxArea)
	if err != nil {
		env.Printf("Error building quad partials: %s\n", err.Error())
		return cover, nil
	}

	if coverOutfile != "" {
		err = CoverDraw(env, cover.Id, coverOutfile)
		if err != nil {
			env.Printf("Error drawing cover: %s\n", err.Error())
			return cover, nil
		}
	}

	return cover, macro
}

func macroQuadBuildPartials(env environment.Environment, cover *model.Cover, macro *model.Macro, img *image.Image, size, minDepth, maxDepth, minArea, maxArea int) error {
	coverPartialService := env.ServiceFactory().MustCoverPartialService()
	quadDistService := env.ServiceFactory().MustQuadDistService()

	var err error

	count, err := coverPartialService.Count(cover)
	if err != nil {
		return err
	}
	current := int(count)

	var total, remain int
	if size > 0 {
		total = macroQuadSplitSize(size)
		remain = total - current
	} else {
		total = -1
		remain = 3
	}

	msgSlice := []string{}
	if size > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("%d splits", size))
	}
	if total > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("%d partials", remain))
	}
	if minDepth > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("min depth %d", minDepth))
	}
	if maxDepth > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("max depth %d", maxDepth))
	}
	if minArea > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("min area %d", minArea))
	}
	if maxArea > 0 {
		msgSlice = append(msgSlice, fmt.Sprintf("max area %d", maxArea))
	}
	env.Println(fmt.Sprintf("Building macro quad with %s...", strings.Join(msgSlice, ", ")))

	bar := pb.StartNew(remain)

	largePartialSatisfied := false

	for {
		var coverPartialQuadView *model.CoverPartialQuadView
		if current == 0 {
			// start with initial values
			coverPartialQuadView = &model.CoverPartialQuadView{
				CoverPartial: &model.CoverPartial{
					CoverId: cover.Id,
					X1:      0,
					Y1:      0,
					X2:      cover.Width,
					Y2:      cover.Height,
				},
				QuadDist: &model.QuadDist{
					Depth: 0,
					Area:  0,
					Dist:  0.0,
				},
			}
			current = 1 // fake the first
		} else {
			// first, check to see if we need minDepth/maxArea
			// remember once we've satisfied this requirement
			if !largePartialSatisfied && (minDepth > 0 || maxArea > 0) {
				coverPartialQuadView, err = quadDistService.GetWorst(macro, minDepth, maxArea)
				if err != nil {
					return err
				}
				if coverPartialQuadView == nil {
					largePartialSatisfied = true
				}
			}

			// if we still don't have a coverPartialQuadView,
			// now we check with maxDepth/minArea
			if coverPartialQuadView == nil {
				coverPartialQuadView, err = quadDistService.GetWorst(macro, maxDepth, minArea)
				if err != nil {
					return err
				}
			}

			if coverPartialQuadView == nil {
				return errors.New("Failed to find worst quad dist")
			}

			err = coverPartialService.Delete(coverPartialQuadView.CoverPartial)
			if err != nil {
				return err
			}
		}

		err = macroQuadSplit(env, macro, coverPartialQuadView, img)
		if err != nil {
			return err
		}
		current += 3

		if total == -1 {
			remain += 3
			bar.Total = int64(remain)
		}
		bar.Set(current)

		if env.Cancel() {
			break
		}

		if current >= remain {
			break
		}
	}

	if env.Cancel() {
		return errors.New("Cancelled")
	}

	bar.Finish()
	return nil
}

func macroQuadCreateCover(env environment.Environment, width, height int) (*model.Cover, error) {
	aspectService := env.ServiceFactory().MustAspectService()
	coverService := env.ServiceFactory().MustCoverService()

	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return nil, err
	}

	cover := &model.Cover{
		AspectId: aspect.Id,
		Width:    width,
		Height:   height,
	}

	err = coverService.Insert(cover)
	if err != nil {
		return nil, err
	}

	return cover, nil
}

func macroQuadSplit(env environment.Environment, macro *model.Macro, coverPartialQuadView *model.CoverPartialQuadView, img *image.Image) error {
	coverPartials, err := macroQuadBuildCoverPartials(env, coverPartialQuadView)
	if err != nil {
		return err
	}

	macroPartials, err := macroQuadBuildMacroPartials(env, macro, coverPartials, img)
	if err != nil {
		return err
	}

	return macroQuadBuildQuadDist(env, coverPartials, macroPartials, coverPartialQuadView.QuadDist, img)
}

func macroQuadBuildCoverPartials(env environment.Environment, coverPartialQuadView *model.CoverPartialQuadView) ([]*model.CoverPartial, error) {
	aspectService := env.ServiceFactory().MustAspectService()
	coverPartialService := env.ServiceFactory().MustCoverPartialService()

	x1 := coverPartialQuadView.CoverPartial.X1
	y1 := coverPartialQuadView.CoverPartial.Y1
	x2 := coverPartialQuadView.CoverPartial.X2
	y2 := coverPartialQuadView.CoverPartial.Y2

	midX := ((x2 - x1) / 2) + x1
	midY := ((y2 - y1) / 2) + y1

	coverPartials := make([]*model.CoverPartial, 4)

	for i, pt := range [][]int{
		[]int{x1, y1, midX, midY},
		[]int{midX, y1, x2, midY},
		[]int{x1, midY, midX, y2},
		[]int{midX, midY, x2, y2},
	} {
		cp := &model.CoverPartial{
			CoverId: coverPartialQuadView.CoverPartial.CoverId,
			X1:      pt[0],
			Y1:      pt[1],
			X2:      pt[2],
			Y2:      pt[3],
		}

		aspect, err := aspectService.FindOrCreate(cp.Width(), cp.Height())
		if err != nil {
			return nil, err
		}

		cp.AspectId = aspect.Id

		err = coverPartialService.Insert(cp)
		if err != nil {
			return nil, err
		}

		coverPartials[i] = cp
	}

	return coverPartials, nil
}

func macroQuadBuildMacroPartials(env environment.Environment, macro *model.Macro, coverPartials []*model.CoverPartial, img *image.Image) ([]*model.MacroPartial, error) {
	macroPartialService := env.ServiceFactory().MustMacroPartialService()

	macroPartials := make([]*model.MacroPartial, 4)
	sem := make(chan bool, 4)

	for idx, coverPartial := range coverPartials {
		sem <- true
		go func(i int, cp *model.CoverPartial) {
			macroPartial := model.MacroPartial{
				MacroId:        macro.Id,
				CoverPartialId: cp.Id,
				AspectId:       cp.AspectId,
			}

			macroPartial.Pixels = util.GetImgPartialLab(img, cp)
			err := macroPartialService.Insert(&macroPartial)
			if err != nil {
				macroPartials[i] = nil
				<-sem
				return
			}

			macroPartials[i] = &macroPartial
			<-sem
		}(idx, coverPartial)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(sem)

	for _, mp := range macroPartials {
		if mp == nil {
			return nil, errors.New("Failed to create macro partial")
		}
	}

	return macroPartials, nil
}

func macroQuadBuildQuadDist(env environment.Environment, coverPartials []*model.CoverPartial, macroPartials []*model.MacroPartial, parent *model.QuadDist, img *image.Image) error {
	quadDistService := env.ServiceFactory().MustQuadDistService()

	sem := make(chan bool, 4)
	errs := false

	for idx := 0; idx < 4; idx++ {
		sem <- true
		go func(i int) {
			quadDist := &model.QuadDist{
				MacroPartialId: macroPartials[i].Id,
				Depth:          parent.Depth + 1,
				Area:           coverPartials[i].Area(),
				Dist:           util.GetImgAvgDist(img, coverPartials[i]),
			}
			err := quadDistService.Insert(quadDist)
			if err != nil {
				errs = true
			}
			<-sem
		}(idx)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(sem)

	if errs {
		return errors.New("Failed to create macro quad dist")
	}

	return nil
}

// macroQuadSplitSize returns the total number of cover partials produced
// from splitting the rectangle n times
func macroQuadSplitSize(n int) int {
	return 4 + 3*n
}

func macroQuadNewMinDepthSplits(minDepth int, cache map[int]int) (int, map[int]int) {
	if minDepth <= 0 {
		return 0, cache
	} else if minDepth == 1 {
		return 1, cache
	} else if splits, ok := cache[minDepth]; ok {
		return splits, cache
	} else {
		rSplits, cache := macroQuadNewMinDepthSplits(minDepth-1, cache)
		splits = rSplits * 4
		cache[minDepth] = splits
		return splits, cache
	}
}

func macroQuadMinDepthSplits(minDepth int) int {
	sum := 0
	cache := make(map[int]int)
	for i := 0; i <= minDepth; i++ {
		var v int
		v, cache = macroQuadNewMinDepthSplits(i, cache)
		sum += v
	}
	return sum
}

// macroQuadFixArgs attempts to produce sane default values from user parameters
// for size, minDepth, maxDepth, minArea, and maxArea
// val == 0 is unrestricted
// val > 0 sets explicitly
// val == -1 (<0) calculates optimal
func macroQuadFixArgs(width, height, size, minDepth, maxDepth, minArea, maxArea int) (int, int, int, int, int, error) {
	var cSize, cMinDepth, cMaxDepth, cMinArea, cMaxArea int

	area := width * height
	normalDim := math.Sqrt(float64(area))

	// cSize is the number of times we split the rectangle
	if size < 0 {
		// set size to 2/5 root of total number of pixels
		cSize = util.Round(math.Pow(float64(area), 0.4))
	} else {
		cSize = size
	}

	if minArea < 0 {
		// min length is the smallest length of a macro partial that we can tolerate
		// it is the bigger of normalDim cut into 85 partials or 35px
		minLength := math.Max(normalDim/float64(85), float64(35))
		cMinArea = util.Round(minLength * minLength)
	} else {
		cMinArea = minArea
	}

	if maxArea < 0 {
		if normalDim <= 1225.0 {
			// we don't restrict maxArea for relatively small macros
			cMaxArea = 0
		} else {
			// max size is the largest length of a macro partial that we can tolerate
			// it is the smaller of aveDim cut into 6 partials, and 600px
			maxSize := util.Round(math.Min(normalDim/float64(6), float64(600)))
			cMaxArea = maxSize * maxSize
		}
	} else {
		cMaxArea = maxArea
	}

	if cMaxArea > 0 && cMinArea > cMaxArea {
		return 0, 0, 0, 0, 0, fmt.Errorf("min-area %d cannot be greater than max-area %d", cMinArea, cMaxArea)
	}

	if minDepth < 0 {
		cMinDepth = 0
		// do not restrict minDepth is size is not restricted
		if cSize > 0 {
			// target a minDepth that produces approx 1/60 the total number of cover partials
			totalPartials := macroQuadSplitSize(cSize)
			minDepthSizeTarget := util.Round(float64(totalPartials) / 60.0)

			// increment cMinDepth until the highest value where it doesn't exceed minDepthSizeTarget
			splits := 0
			for depth := 0; splits <= minDepthSizeTarget; depth += 1 {
				splits = macroQuadMinDepthSplits(depth)
				if splits <= minDepthSizeTarget {
					cMinDepth = depth
				}
			}
		}
	} else {
		cMinDepth = minDepth
	}

	if cMinDepth > 0 {
		var splits int
		splits = macroQuadMinDepthSplits(cMinDepth)
		if splits > cSize {
			return 0, 0, 0, 0, 0, fmt.Errorf("min-depth %d too large for size %d", cMinDepth, cSize)
		}
	}

	if maxDepth < 0 {
		// we want a maxDepth such that
		// minArea^(1/2) = normalDim / depth^2
		// solve for depth
		var minLength float64
		if cMinArea > 0 {
			minLength = math.Sqrt(float64(cMinArea))
		} else {
			// see minArea calculation above
			minLength = math.Max(normalDim/float64(85), float64(35))
		}
		cMaxDepth = util.Round(math.Sqrt(normalDim / minLength))

		if cMinDepth > 0 && cMinDepth >= cMaxDepth {
			// fall back to basing this off cMinDepth
			cMaxDepth = util.Round(float64(cMinDepth) * 1.2)
			if cMinDepth == cMaxDepth {
				cMaxDepth += 1
			}
		}
	} else {
		cMaxDepth = maxDepth
	}

	if cMaxDepth > 0 && cMinDepth > cMaxDepth {
		return 0, 0, 0, 0, 0, fmt.Errorf("min-depth %d cannot be greater than max-depth %d", cMinDepth, cMaxDepth)
	}

	return cSize, cMinDepth, cMaxDepth, cMinArea, cMaxArea, nil
}
