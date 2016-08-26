package controller

//func MacroQuad(env environment.Environment, path string, coverWidth, coverHeight, num int, outfile string) (*model.Cover, *model.Macro) {
//	var myCoverWidth, myCoverHeight int
//
//	if coverWidth < 0 || coverHeight < 0 {
//		env.Println("Cover width and height must not be less than zero")
//		return nil, nil
//	}
//
//	if coverWidth > 0 && coverHeight > 0 {
//		myCoverWidth = coverWidth
//		myCoverHeight = coverHeight
//	} else {
//		aspectService, err := env.AspectService()
//		if err != nil {
//			env.Printf("Error getting aspect service: %s\n", err.Error())
//			return nil, nil
//		}
//
//		aspect, err := macroAspectGetImageAspect(path, aspectService)
//		if err != nil {
//			env.Printf("Error getting aspect: %s\n", err.Error())
//			return nil, nil
//		}
//
//		if coverWidth == 0 {
//			myCoverWidth = aspect.RoundWidth(coverHeight)
//		} else {
//			myCoverWidth = coverWidth
//		}
//
//		if coverHeight == 0 {
//			myCoverHeight = aspect.RoundHeight(coverWidth)
//		} else {
//			myCoverHeight = coverHeight
//		}
//	}
//
//	cover := CoverAspect(env, myCoverWidth, myCoverHeight, partialWidth, partialHeight, num)
//	if cover == nil {
//		env.Println("Failed to create cover")
//		return nil, nil
//	}
//	macro := Macro(env, path, cover.Id, outfile)
//	if macro == nil {
//		env.Println("Failed to create macro")
//		return cover, nil
//	}
//
//	return cover, macro
//}
