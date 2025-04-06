package game

type MovableSolid struct{}

func (MovableSolid) Create(cell *Cell) error {
	return nil
}

func (MovableSolid) Update(cell *Cell) error {
	bottom, err := cell.GetCell(0, 1)
	if err != nil {
		return nil
	}

	if bottom.Data().Role == ROLE_AIR {
		bottom.Type, cell.Type = cell.Type, bottom.Type
	}

	return nil
}
