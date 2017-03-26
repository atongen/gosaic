package controller

import (
	"fmt"
	"github.com/atongen/gosaic/environment"
)

func MosaicList(env environment.Environment) {
	mosaicService := env.ServiceFactory().MustMosaicService()
	mosaicPartialService := env.ServiceFactory().MustMosaicPartialService()

	mosaics, err := mosaicService.FindAll("mosaics.id desc")
	if err != nil {
		env.Printf("Error finding mosaics: %s\n", err.Error())
		return
	}
	if len(mosaics) == 0 {
		// we are done
		return
	}

	for _, mosaic := range mosaics {
		fmt.Println(mosaic)
		num, err := mosaicPartialService.Count(mosaic)
		if err != nil {
			env.Printf("Error counting mosaic partials: %s\n", err.Error())
			return
		}
		fmt.Printf("\tNum partials: %d\n", num)
	}
}
