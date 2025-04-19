package game

import "math/rand"

type Liquid struct{}

func (Liquid) IsA(kind string) bool {
	return kind == "Dynamic" || kind == "Liquid"
}

func (Liquid) Create(cell *Cell) error {
	return nil
}

func (Liquid) Update(cell *Cell) error {
	bottom, err := cell.GetCell(0, 1)
	if err != nil {
		return nil
	}

	if cell.CanFallInto(bottom) {
		cell.Switch(bottom)
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
		if cell.CanMoveInto(bottomLeft) {
			cell.Switch(bottomLeft)
			return nil
		} else {
			dir *= -1
		}
	}

	bottomRight, err := cell.GetCell(dir, 0)
	if err != nil {
		dir *= -1
	} else {
		if cell.CanMoveInto(bottomRight) {
			cell.Switch(bottomRight)
			return nil
		} else {
			dir *= -1
		}
	}

	return nil
}
