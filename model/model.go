package model

import (
  "database/sql"
)

type Model interface {
  GetId() uint
  Create(db *sql.DB) bool
  Update(db *sql.DB) bool
  Delete(db *sql.DB) bool
}
