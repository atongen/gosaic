package model

type MacroGidxView struct {
	MacroPartial *MacroPartial
	GidxPartial  *GidxPartial
}

func (macroGidxView *MacroGidxView) PartialComparison() (*PartialComparison, error) {
	err := macroGidxView.MacroPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	err = macroGidxView.GidxPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	dist, err := PixelDist(macroGidxView.MacroPartial, macroGidxView.GidxPartial)
	if err != nil {
		return nil, err
	}

	return &PartialComparison{
		MacroPartialId: macroGidxView.MacroPartial.Id,
		GidxPartialId:  macroGidxView.GidxPartial.Id,
		Dist:           dist,
	}, nil
}
