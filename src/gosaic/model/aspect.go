package model

type Aspect struct {
	Id      int64 `db:"id"`
	Columns int   `db:"columns"`
	Rows    int   `db:"rows"`
}

func NewAspect(columns int, rows int) *Aspect {
	return &Aspect{
		Columns: columns,
		Rows:    rows,
	}
}

func (a *Aspect) SetAspect(width int, height int) (int, int) {
	c, r := calculateAspect(width, height)
	a.Columns = c
	a.Rows = r
	return a.Columns, a.Rows
}
