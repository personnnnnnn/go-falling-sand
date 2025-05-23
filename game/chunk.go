package game

import (
	"fmt"
	"go-falling-sand/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Chunk struct {
	X, Y      int
	Game      *Game
	Cells     []Cell
	CellOrder []int
}

func NewChunk(game *Game, x, y int) Chunk {
	chunk := Chunk{}

	chunk.Game = game

	chunk.X = x
	chunk.Y = y

	chunk.Cells = make([]Cell, game.ChunkArea())
	chunk.CellOrder = make([]int, game.ChunkArea())

	for x := range game.ChunkWidth {
		for y := range game.ChunkHeight {

			i := game.CalculateCellIndex(x, y)

			chunk.CellOrder[i] = i

			var cellType int

			worldX := x + chunk.X*game.ChunkWidth
			worldY := y + chunk.Y*game.ChunkHeight

			if worldX == 0 || worldY == 0 || worldX == game.TotalWidth()-1 || worldY == game.TotalHeight()-1 {
				cellType = game.WallElement
			} else {
				cellType = game.AirElement
			}

			cell := Cell{
				X: x, Y: y,
				Type:  cellType,
				Chunk: &chunk,
			}

			chunk.Cells[i] = cell
		}
	}

	return chunk
}

func (chunk *Chunk) GetCell(cellX, cellY int) (*Cell, error) {
	if cellX < 0 || cellY < 0 || cellX >= chunk.Game.ChunkWidth || cellY >= chunk.Game.ChunkHeight {
		return nil, fmt.Errorf("there is no cell in chunk at local position %v %v", cellX, cellY)
	}
	return &chunk.Cells[chunk.Game.CalculateCellIndex(cellX, cellY)], nil
}

func (chunk *Chunk) Update() error {
	util.Shuffle(chunk.CellOrder)
	for i := range chunk.CellOrder {
		i = chunk.CellOrder[i]
		cell := &chunk.Cells[i]
		if err := cell.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (chunk *Chunk) Draw(screen *ebiten.Image) {
	for x := range chunk.Game.ChunkWidth {
		for y := range chunk.Game.ChunkHeight {
			i := chunk.Game.CalculateCellIndex(x, y)

			cell := chunk.Cells[i]

			vector.DrawFilledRect(
				screen,
				float32(x+chunk.X*chunk.Game.ChunkWidth)*chunk.Game.CellSize+chunk.Game.SideBarLength,
				float32(y+chunk.Y*chunk.Game.ChunkHeight)*chunk.Game.CellSize,
				chunk.Game.CellSize,
				chunk.Game.CellSize,
				chunk.Game.ElementData[cell.Type].Color,
				false,
			)
		}
	}
}
