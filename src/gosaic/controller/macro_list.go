package controller

import "gosaic/environment"

func MacroList(env environment.Environment) {
	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
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
		env.Println(macro)
	}
}
