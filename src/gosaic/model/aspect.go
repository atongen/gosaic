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
	var d int = gcd(width, height)
	a.Columns = width / d
	a.Rows = height / d
	return a.Columns, a.Rows
}

func AspectsToInterface(aspects []*Aspect) []interface{} {
	n := len(aspects)
	interfaces := make([]interface{}, n)
	for i := 0; i < n; i++ {
		interfaces[i] = interface{}(aspects[i])
	}
	return interfaces
}
