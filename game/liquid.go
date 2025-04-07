package game

import "math/rand"

type Liquid struct{}

func (Liquid) Create(cell *Cell) error {
	return nil
}

func (Liquid) Update(cell *Cell) error {
	bottom, err := cell.GetCell(0, 1)
	if err != nil {
		return nil
	}

	if bottom.ElementData().Role == ROLE_AIR {
		bottom.Type, cell.Type = cell.Type, bottom.Type
		bottom.Data, cell.Data = cell.Data, bottom.Data
		return nil
	}

	dir := 1
	if rand.Float32() < 0.5 {
		dir *= -1
	}

	bottomLeft, err := cell.GetCell(dir, 0)
	if err != nil {
		dir *= -1
	} else {
		if bottomLeft.ElementData().Role == ROLE_AIR {
			bottomLeft.Type, cell.Type = cell.Type, bottomLeft.Type
			bottomLeft.Data, cell.Data = cell.Data, bottomLeft.Data
			return nil
		} else {
			dir *= -1
		}
	}

	bottomRight, err := cell.GetCell(dir, 0)
	if err != nil {
		dir *= -1
	} else {
		if bottomRight.ElementData().Role == ROLE_AIR {
			bottomRight.Type, cell.Type = cell.Type, bottomRight.Type
			bottomRight.Data, cell.Data = cell.Data, bottomRight.Data
			return nil
		} else {
			dir *= -1
		}
	}

	return nil
}
