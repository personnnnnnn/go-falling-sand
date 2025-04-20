package game

import (
	"go-falling-sand/util"
	"math/rand"
)

const CUSTOM_DO_NOTHING = 0
const CUSTOM_END = 1

type Condition interface {
	Satisfied(cell *Cell) (bool, error)
}

type Action interface {
	Act(cell *Cell) (int, error)
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
		if res, err := result.Act(cell); err != nil {
			return err
		} else if res == CUSTOM_END {
			return nil
		}
	}

	return nil
}

type TurnInto struct {
	ID int
}

func (kind *TurnInto) Act(cell *Cell) (int, error) {
	cell.Type = kind.ID
	return CUSTOM_DO_NOTHING, nil
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
			if !(x == 0 && y == 0) {
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

type DirectlyTouching struct {
	ID int
}

func (kind *DirectlyTouching) Satisfied(cell *Cell) (bool, error) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 || y == 0 {
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

func (kind *Emit) Act(cell *Cell) (int, error) {
	dx, dy := util.GetRandomDir()
	other, err := cell.GetCell(dx, dy)
	if err != nil {
		return CUSTOM_DO_NOTHING, nil
	}
	if other.ElementData().Role != ROLE_AIR {
		return CUSTOM_DO_NOTHING, nil
	}
	other.Type = kind.ID
	return CUSTOM_DO_NOTHING, nil
}

type End struct{}

func (End) Act(cell *Cell) (int, error) {
	return CUSTOM_END, nil
}
