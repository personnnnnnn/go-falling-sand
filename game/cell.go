package game

type Cell struct {
	X, Y  int
	Type  int
	Chunk *Chunk
	Data  *[]int
}

func (cell *Cell) Update() error {
	kind := cell.Game().ElementData[cell.Type].Kind
	if kind == nil {
		return nil
	}
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

func (cell *Cell) ElementData() *ElementData {
	data := cell.Game().ElementData[cell.Type]
	return &data
}
