package game

type Dust struct {
	Weight float32
}

func (Dust) IsA(kind string) bool {
	return kind == "MovableSolid" || kind == "Dust"
}

func (Dust) Create(cell *Cell) error {
	return nil
}

func (dust *Dust) Update(cell *Cell) error {
	bottom, err := cell.GetCell(0, 1)
	if err == nil && !(cell.CanFallInto(bottom) && !bottom.IsSolid()) {
		return MovableSolid{}.Update(cell)
	} else {
		top, err := cell.GetCell(0, -1)
		if err == nil && !(cell.CanFallInto(top) && !top.IsSolid()) {
			return MovableSolid{}.Update(cell)
		}
	}

	return (&Gas{dust.Weight}).Update(cell)
}
