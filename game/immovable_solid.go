package game

type ImmovableSolid struct{}

func (ImmovableSolid) IsA(kind string) bool {
	return kind == "Solid" || kind == "ImmovableSolid"
}

func (ImmovableSolid) Create(cell *Cell) error {
	return nil
}

func (ImmovableSolid) Update(cell *Cell) error {
	return nil
}
