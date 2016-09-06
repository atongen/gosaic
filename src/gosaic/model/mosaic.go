package model

import (
	"fmt"
	"time"
)

type Mosaic struct {
	Id         int64     `db:"id"`
	Name       string    `db:"name"`
	MacroId    int64     `db:"macro_id"`
	IsComplete bool      `db:"is_complete"`
	CreatedAt  time.Time `db:"created_at"`
}

func (m *Mosaic) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, MacroId: %d, IsComplete: %v, CreatedAt: %s",
		m.Id, m.Name, m.MacroId, m.IsComplete, m.CreatedAt)
}
