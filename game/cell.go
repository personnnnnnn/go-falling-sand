package game

import "errors"

type Cell struct {
	X, Y        int
	Type        int
	Chunk       *Chunk
	Data        *[]int
	UpdateCycle bool
}

func (cell *Cell) HasUpdated() bool {
	return cell.UpdateCycle != cell.Game().UpdateCycle
}

func (cell *Cell) Update() error {
	if cell.HasUpdated() {
		return nil
	}
	cell.UpdateCycle = !cell.UpdateCycle
	kind := cell.ElementData().Kind
	return kind.Update(cell)
}

func (cell *Cell) GetCell(relativeX, relativeY int) (*Cell, error) {
	return cell.Game().GetCell(cell.WorldX()+relativeX, cell.WorldY()+relativeY)
}

func (cell *Cell) Game() *Game {
	return cell.Chunk.Game
}

func (cell *Cell) WorldX() int {
	return cell.X + cell.Chunk.X*cell.Game().ChunkWidth
}

func (cell *Cell) WorldY() int {
	return cell.Y + cell.Chunk.Y*cell.Game().ChunkHeight
}

func (cell *Cell) ElementData() ElementData {
	return cell.Game().ElementData[cell.Type]
}

func (cell *Cell) Index() int {
	return cell.Game().CalculateCellIndex(cell.X, cell.Y)
}

func (cell *Cell) Switch(other *Cell) (*Cell, error) {
	if other == nil {
		return nil, errors.New("can't switch places with nil cell")
	}

	cell.Type, other.Type = other.Type, cell.Type
	cell.UpdateCycle, other.UpdateCycle = other.UpdateCycle, cell.UpdateCycle

	if !cell.HasUpdated() {
		err := cell.Update()
		if err != nil {
			return nil, err
		}
	}

	return other, nil
}

func (cell *Cell) CanFallInto(other *Cell) bool {
	return other.ElementData().Bouyancy < cell.ElementData().Bouyancy
}

func (cell *Cell) CanMoveInto(other *Cell) bool {
	return cell.CanFallInto(other) && !other.HasUpdated()
}

func (cell *Cell) IsA(kind string) bool {
	return cell.ElementData().Kind.IsA(kind)
}

func (cell *Cell) IsSolid() bool {
	return cell.IsA("Solid")
}
