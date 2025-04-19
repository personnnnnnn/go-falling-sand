package game

import "math/rand"

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
