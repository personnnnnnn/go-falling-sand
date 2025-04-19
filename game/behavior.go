package game

type ElementKind interface {
	Create(cell *Cell) error
	Update(cell *Cell) error
	IsA(kind string) bool
}
