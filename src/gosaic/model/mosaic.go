package model

import (
	"fmt"
	"time"
)

type Mosaic struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	MacroId   int64     `db:"macro_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (m *Mosaic) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, MacroId: %d, CreatedAt: %s",
		m.Id, m.Name, m.MacroId, m.CreatedAt)
}
