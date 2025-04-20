package game

import (
	"go-falling-sand/util"
	"math/rand"
)

type Condition interface {
	Satisfied(cell *Cell) (bool, error)
}

type Action interface {
	Act(cell *Cell) error
}

type Reaction struct {
	Conditions []Condition
	Actions    []Action
}

func (*Reaction) IsA(kind string) bool {
	return kind == "CustomKind"
}

func (kind *Reaction) Create(cell *Cell) error {
	return nil
}

func (kind *Reaction) Update(cell *Cell) error {
	for i := range kind.Conditions {
		condition := kind.Conditions[i]
		res, err := condition.Satisfied(cell)
		if err != nil {
			return err
		}
		if !res {
			return nil
		}
	}

	for i := range kind.Actions {
		result := kind.Actions[i]
		if err := result.Act(cell); err != nil {
			return err
		}
	}

	return nil
}

type TurnInto struct {
	ID int
}

func (kind *TurnInto) Act(cell *Cell) error {
	cell.Type = kind.ID
	return nil
}

type Chance struct {
	Chance float32
}

func (kind *Chance) Satisfied(cell *Cell) (bool, error) {
	return rand.Float32() < kind.Chance, nil
}

type Touching struct {
	ID int
}

func (kind *Touching) Satisfied(cell *Cell) (bool, error) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x != 0 && y != 0 {
				if other, err := cell.GetCell(x, y); err != nil {
					continue
				} else if other.Type == kind.ID {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

type Emit struct {
	ID int
}

func (kind *Emit) Act(cell *Cell) error {
	dx, dy := util.GetRandomDir()
	other, err := cell.GetCell(dx, dy)
	if err != nil {
		return nil
	}
	if other.ElementData().Role != ROLE_AIR {
		return nil
	}
	other.Type = kind.ID
	return nil
}
