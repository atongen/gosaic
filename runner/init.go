package runner

import (
	"fmt"
	"github.com/atongen/gosaic"
	"github.com/atongen/gosaic/database"
)

type Init Run

func (init Init) Execute() error {
	isProj, err := gosaic.IsProject(init.Path)
	if err != nil {
		return err
	} else if isProj {
		fmt.Printf("Existing project: %s\n", init.Path)
	} else {
		fmt.Printf("New project: %s\n", init.Path)
	}

	// ensure directory is present
	err = gosaic.MkProjectDir(init.Path)
	if err != nil {
		return err
	}

	// migrate the database
	err = database.Migrate(gosaic.DbPath(init.Path))
	if err != nil {
		return err
	}

	return nil
}
