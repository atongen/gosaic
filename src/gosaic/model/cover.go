package model

import "fmt"

type Cover struct {
	Id     int64  `db:"id"`
	Type   string `db:"type"`
	Name   string `db:"name"`
	Width  uint   `db:"width"`
	Height uint   `db:"height"`
}

func (c *Cover) String() string {
	return fmt.Sprintf("%s (%s) %dx%d", c.Name, c.Type, c.Width, c.Height)
}
