package runner

import (
	"fmt"
	"os"
)

type Status Run

func (status Status) Execute() error {
	dbPath := status.Project.DbPath()
	_, err := os.Stat(dbPath)
	if err == nil {
		fmt.Printf("Database exists: %s\n", dbPath)
	} else {
		return err
	}

	return nil
}
