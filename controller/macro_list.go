package controller

import (
	"fmt"
	"github.com/atongen/gosaic/environment"
)

func MacroList(env environment.Environment) {
	macroService := env.ServiceFactory().MustMacroService()
	macroPartialService := env.ServiceFactory().MustMacroPartialService()
	coverService := env.ServiceFactory().MustCoverService()
	aspectService := env.ServiceFactory().MustAspectService()

	macros, err := macroService.FindAll("macros.id desc")
	if err != nil {
		env.Printf("Error finding macros: %s\n", err.Error())
		return
	}
	if len(macros) == 0 {
		// we are done
		return
	}

	for _, macro := range macros {
		fmt.Printf("ID: %d\n\tPath: %s (%dx%d)\n", macro.Id, macro.Path, macro.Width, macro.Height)

		aspectIds, err := macroPartialService.AspectIds(macro.Id)
		if err != nil {
			fmt.Printf("\tError getting aspect ids for macro partials: %s\n", err.Error())
		} else {
			aspects, err := aspectService.FindIn(aspectIds)
			if err != nil {
				fmt.Printf("\tError getting aspects for macro partials: %s\n", err.Error())
			} else {
				fmt.Printf("\tPartial aspects: ")
				for i, aspect := range aspects {
					fmt.Printf("%dx%d", aspect.Columns, aspect.Rows)
					if i < len(aspects)-1 {
						fmt.Printf(", ")
					}
				}
				fmt.Printf("\n")
			}
		}

		cover, err := coverService.Get(macro.CoverId)
		if err != nil {
			fmt.Printf("\tError getting cover details: %s\n", err.Error())
		} else {
			fmt.Printf("\tCover: %s\n", cover)
		}
	}
}
