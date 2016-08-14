package controller

import (
	"fmt"
	"gosaic/environment"
)

func MacroList(env environment.Environment) {
	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
		return
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error getting macro partial service: %s\n", err.Error())
		return
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error getting cover service: %s\n", err.Error())
		return
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return
	}

	macros, err := macroService.FindAll("macros.name ASC")
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
