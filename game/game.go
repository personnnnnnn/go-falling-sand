package game

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"go-falling-sand/xml_handler"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var Dimensions = struct {
	Width, Height, ScreenWidth, ScreenHeight int
}{400, 400, 600, 600}

func HexStringToColor(s string) (color.Color, error) {
	if len(s) != 7 || s[0] != '#' {
		return color.Black, errors.New("invalid hex color")
	}

	rStr := s[1:3]
	gStr := s[3:5]
	bStr := s[5:7]

	r, err := strconv.ParseUint(rStr, 16, 8)
	if err != nil {
		return color.Black, fmt.Errorf("invalid red component: %v", err)
	}
	g, err := strconv.ParseUint(gStr, 16, 8)
	if err != nil {
		return color.Black, fmt.Errorf("invalid green component: %v", err)
	}
	b, err := strconv.ParseUint(bStr, 16, 8)
	if err != nil {
		return color.Black, fmt.Errorf("invalid blue component: %v", err)
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

func StringToColor(s string) (color.Color, error) {
	if len(s) == 0 {
		return color.Black, errors.New("can't give an empty string")
	}

	if s[0] == '#' {
		return HexStringToColor(s)
	}

	if col, ok := ColorMap[s]; ok {
		return col, nil
	}

	return color.Black, fmt.Errorf("invalid color name: '%v'", s)
}

const ROLE_WALL = "wall"
const ROLE_AIR = "air"
const ROLE_NONE = "none"

type ElementData struct {
	Color           color.Color
	Name            string
	ElementTypeName string
	ElementTypeID   int
	Role            string
}

type Game struct {
	elementIdCounter        int
	AirElement              int
	WallElement             int
	ElementTypes            map[string]int
	ElementData             map[int]ElementData
	Width, Height           int
	ChunkWidth, ChunkHeight int
	Chunks                  []Chunk
	CellSize                float32
}

type Chunk struct {
	X, Y  int
	Game  *Game
	Cells []Cell
}

func (g *Game) TotalWidth() int {
	return g.Width * g.ChunkWidth
}

func (g *Game) TotalHeight() int {
	return g.Height * g.ChunkHeight
}

func (g *Game) DefineElement(elementTypeName string, color string, name string, role string) error {
	var col, colErr = StringToColor(color)
	if colErr != nil {
		return colErr
	}

	if role != ROLE_AIR && role != ROLE_WALL && role != ROLE_NONE {
		role = ROLE_NONE
	}

	g.ElementTypes[elementTypeName] = g.elementIdCounter
	g.ElementData[g.elementIdCounter] = ElementData{
		Color:           col,
		Name:            name,
		ElementTypeName: elementTypeName,
		ElementTypeID:   g.elementIdCounter,
		Role:            role,
	}

	if role == ROLE_AIR {
		g.AirElement = g.elementIdCounter
	} else if role == ROLE_WALL {
		g.WallElement = g.elementIdCounter
	}

	g.elementIdCounter++

	return nil
}

type Cell struct {
	X, Y  int
	Type  int
	Chunk *Chunk
}

func (game *Game) ChunkArea() int {
	return game.ChunkWidth * game.ChunkHeight
}

func (game *Game) WorldArea() int {
	return game.Width * game.Height
}

func NewGame(width, height int, chunkWidth, chunkHeight int, cellSize float32, xmlData []byte) (*Game, error) {
	game := &Game{}

	game.elementIdCounter = 0

	game.Width = width
	game.Height = height

	game.ChunkWidth = chunkWidth
	game.ChunkHeight = chunkHeight
	game.CellSize = cellSize

	game.ElementData = map[int]ElementData{}
	game.ElementTypes = map[string]int{}

	var commands xmlhandler.XMLElementList
	xml.Unmarshal(xmlData, &commands)

	for i := range commands.Elements {
		command := commands.Elements[i]

		display := command.Display
		if display == nil {
			display = &xmlhandler.XMLDisplay{}
		}

		col := display.Color
		if col == "" {
			col = "white"
		}

		name := display.Name
		if name == "" {
			name = command.Name
		}

		if err := game.DefineElement(command.Name, col, name, command.Role); err != nil {
			return nil, err
		}
	}

	game.Chunks = make([]Chunk, game.WorldArea())
	for x := 0; x < game.Width; x++ {
		for y := 0; y < game.Height; y++ {
			i := game.CalculateChunkIndex(x, y)

			game.Chunks[i] = NewChunk(game, x, y)
		}
	}

	return game, nil
}

func NewChunk(game *Game, x, y int) Chunk {
	chunk := Chunk{}

	chunk.Game = game

	chunk.X = x
	chunk.Y = y

	chunk.Cells = make([]Cell, game.ChunkArea())

	for x := 0; x < game.ChunkWidth; x++ {
		for y := 0; y < game.ChunkHeight; y++ {

			i := game.CalculateCellIndex(x, y)

			var cellType int

			worldX := x + chunk.X*game.ChunkWidth
			worldY := y + chunk.Y*game.ChunkHeight

			if worldX == 0 || worldY == 0 || worldX == game.TotalWidth()-1 || worldY == game.TotalHeight()-1 {
				cellType = game.WallElement
			} else {
				cellType = game.AirElement
			}

			chunk.Cells[i] = Cell{
				X: x, Y: y,
				Type:  cellType,
				Chunk: &chunk,
			}
		}
	}

	return chunk
}

func (g *Game) CalculateChunkIndex(x, y int) int {
	return x + y*g.Width
}

func (g *Game) CalculateCellIndex(x, y int) int {
	return x + y*g.ChunkWidth
}

func (game *Game) Layout(outsizeWidth, outsizeHeight int) (int, int) {
	return Dimensions.Width, Dimensions.Height
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for x := 0; x < game.Width; x++ {
		for y := 0; y < game.Height; y++ {
			i := game.CalculateChunkIndex(x, y)

			chunk := game.Chunks[i]
			chunk.Draw(screen)
		}
	}
}

func (chunk *Chunk) Draw(screen *ebiten.Image) {
	for x := 0; x < chunk.Game.ChunkWidth; x++ {
		for y := 0; y < chunk.Game.ChunkHeight; y++ {
			i := chunk.Game.CalculateCellIndex(x, y)

			cell := chunk.Cells[i]

			vector.DrawFilledRect(
				screen,
				float32(x+chunk.X*chunk.Game.ChunkWidth)*chunk.Game.CellSize,
				float32(y+chunk.Y*chunk.Game.ChunkHeight)*chunk.Game.CellSize,
				chunk.Game.CellSize,
				chunk.Game.CellSize,
				chunk.Game.ElementData[cell.Type].Color,
				false,
			)
		}
	}
}

func (game *Game) Update() error {
	return nil
}
