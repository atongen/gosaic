package model

type MacroGidxView struct {
	MacroPartialId   int64  `db:"macro_partial_id"`
	MacroId          int64  `db:"macro_id"`
	CoverPartialId   int64  `db:"cover_partial_id"`
	AspectId         int64  `db:"aspect_id"`
	MacroPartialData []byte `db:"macro_partial_data"`
	GidxPartialId    int64  `db:"gidx_partial_id"`
	GidxId           int64  `db:"gidx_id"`
	GidxPartialData  []byte `db:"gidx_partial_data"`
}

func (mgv *MacroGidxView) MacroPartial() *MacroPartial {
	return &MacroPartial{
		Id:             mgv.MacroPartialId,
		MacroId:        mgv.MacroId,
		CoverPartialId: mgv.CoverPartialId,
		AspectId:       mgv.AspectId,
		Data:           mgv.MacroPartialData,
	}
}

func (mgv *MacroGidxView) GidxPartial() *GidxPartial {
	return &GidxPartial{
		Id:       mgv.GidxPartialId,
		GidxId:   mgv.GidxId,
		AspectId: mgv.AspectId,
		Data:     mgv.GidxPartialData,
	}
}
