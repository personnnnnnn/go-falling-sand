package game

import (
	"go-falling-sand/util"
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
		dx, dy = util.GetRandomDir()
	} else {
		dx = rand.Intn(3) - 1
		dy = 1
	}

	other, err := cell.GetCell(dx, dy)
	if err == nil {
		if other.IsA("Gas") {
			cell.Switch(other)
			return nil
		}
	}

	return nil
}
