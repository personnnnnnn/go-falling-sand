package game

import "math/rand"

type MovableSolid struct{}

func (MovableSolid) IsA(kind string) bool {
	return kind == "Solid" || kind == "MovableSolid"
}

func (MovableSolid) Create(cell *Cell) error {
	return nil
}

func (MovableSolid) Update(cell *Cell) error {
	bottom, err := cell.GetCell(0, 1)
	if err != nil {
		return nil
	}

	if cell.CanFallInto(bottom) && !bottom.IsSolid() {
		cell.Switch(bottom)
		return nil
	}

	dir := 1
	if rand.Float32() < 0.5 {
		dir *= -1
	}

	bottomLeft, err := cell.GetCell(dir, 1)
	if err != nil {
		dir *= -1
	} else {
		if cell.CanMoveInto(bottomLeft) && !bottomLeft.IsSolid() {
			cell.Switch(bottomLeft)
			return nil
		} else {
			dir *= -1
		}
	}

	bottomRight, err := cell.GetCell(dir, 1)
	if err != nil {
		dir *= -1
	} else {
		if cell.CanMoveInto(bottomRight) && !bottomRight.IsSolid() {
			cell.Switch(bottomRight)
			return nil
		} else {
			dir *= -1
		}
	}

	return nil
}
