package model

import (
  "database/sql"
)

type Model interface {
  GetId() uint
  Create(db *sql.DB) error
  Update(db *sql.DB) error
  Delete(db *sql.DB) error
}
