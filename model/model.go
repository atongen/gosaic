package model

import (
  "fmt"
)

type Model interface (
  TableName() string
)
