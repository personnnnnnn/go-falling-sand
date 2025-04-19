package game

type DefaultKind struct{}

func (DefaultKind) Create(cell *Cell) error {
	return nil
}

func (DefaultKind) Update(cell *Cell) error {
	return nil
}

func (DefaultKind) IsA(kind string) bool {
	return false
}
