package game

import (
	"math/rand"
)

type Gas struct {
	Weight float32
}

func (Gas) IsA(kind string) bool {
	return kind == "Dynamic" || kind == "Gas"
}

func (Gas) Create(cell *Cell) error {
	return nil
}

func (gas *Gas) Update(cell *Cell) error {
	var dx int
	var dy int
	if rand.Float32() > gas.Weight {
		n := rand.Intn(8)
		if n == 0 {
			dx = 0
			dy = 1
		} else if n == 1 {
			dx = 1
			dy = 1
		} else if n == 2 {
			dx = 1
			dy = 0
		} else if n == 3 {
			dx = 1
			dy = -1
		} else if n == 4 {
			dx = 0
			dy = -1
		} else if n == 5 {
			dx = -1
			dy = -1
		} else if n == 6 {
			dx = -1
			dy = 0
		} else {
			dx = -1
			dy = 1
		}
	} else {
		dx = rand.Intn(3) - 1
		dy = 1
	}

	other, err := cell.GetCell(dx, dy)
	if err == nil {
		if cell.CanMoveInto(other) {
			cell.Switch(other)
			return nil
		}
	}

	return nil
}
