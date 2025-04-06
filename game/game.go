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
}{600, 400, 900, 600}

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
	Kind            ElementKind
}

type Game struct {
	elementIdCounter        int
	SideBarLength           float32
	AirElement              int
	WallElement             int
	SelectedElement         int
	ElementTypes            map[string]int
	ElementData             map[int]ElementData
	Width, Height           int
	ChunkWidth, ChunkHeight int
	Chunks                  []Chunk
	CellSize                float32
	ElementScrollBar        ScrollBar
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

func (g *Game) DefineElement(elementTypeName string, colorString string, name string, role string, selectable bool) error {
	index := g.elementIdCounter
	g.elementIdCounter++

	var col, colErr = StringToColor(colorString)
	if colErr != nil {
		return colErr
	}

	if role != ROLE_AIR && role != ROLE_WALL && role != ROLE_NONE {
		role = ROLE_NONE
	}

	g.ElementTypes[elementTypeName] = index
	g.ElementData[index] = ElementData{
		Color:           col,
		Name:            name,
		ElementTypeName: elementTypeName,
		ElementTypeID:   index,
		Role:            role,
	}

	if role == ROLE_AIR {
		g.AirElement = index
	} else if role == ROLE_WALL {
		g.WallElement = index
	}

	if selectable {
		g.ElementScrollBar.AddItem(ScrollBarItem{
			Box: &ScrollBarBox{
				Border:     color.White,
				Inner:      col,
				BorderSize: 3,
			},
			InnerPadding: 3,
			TextColor:    color.White,
			Text:         name,
			Clicked: func(_ *ScrollBarItem, _ int) error {
				g.SelectedElement = index
				return nil
			},
			BeforeDraw: func(item *ScrollBarItem, i int) {
				if g.SelectedElement == index {
					item.Background = color.RGBA{200, 200, 200, 255}
				} else if g.ElementScrollBar.GetHovered() == i {
					item.Background = color.RGBA{150, 150, 150, 255}
				} else {
					item.Background = color.Transparent
				}
			},
		})
	}

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

func NewGame(width, height int, chunkWidth, chunkHeight int, cellSize float32, sideBarLength float32, xmlData []byte) (*Game, error) {
	game := &Game{}

	game.elementIdCounter = 0

	game.SelectedElement = -1 // No item selected

	game.Width = width
	game.Height = height

	game.ChunkWidth = chunkWidth
	game.ChunkHeight = chunkHeight

	game.CellSize = cellSize
	game.SideBarLength = sideBarLength

	game.ElementData = map[int]ElementData{}
	game.ElementTypes = map[string]int{}

	var commands xmlhandler.XMLElementList
	xml.Unmarshal(xmlData, &commands)

	game.ElementScrollBar = NewScrollBar(
		0,
		sideBarLength,
		30,
		10,
		color.RGBA{100, 100, 100, 255},
		len(commands.Elements),
	)

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

		if err := game.DefineElement(command.Name, col, name, command.Role, display.Selectable); err != nil {
			return nil, err
		}
	}

	game.Chunks = make([]Chunk, game.WorldArea())
	for x := range game.Width {
		for y := range game.Height {
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

	for x := range game.ChunkWidth {
		for y := range game.ChunkHeight {

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
	screen.Fill(color.Gray{100})

	for x := range game.Width {
		for y := range game.Height {
			i := game.CalculateChunkIndex(x, y)

			chunk := &game.Chunks[i]
			chunk.Draw(screen)
		}
	}

	game.ElementScrollBar.Draw(screen)
}

func (chunk *Chunk) Draw(screen *ebiten.Image) {
	for x := range chunk.Game.ChunkWidth {
		for y := range chunk.Game.ChunkHeight {
			i := chunk.Game.CalculateCellIndex(x, y)

			cell := &chunk.Cells[i]

			col := chunk.Game.ElementData[cell.Type].Color
			if c, err := chunk.Game.GetHoveredCell(); err == nil && c == cell {
				col = color.RGBA{255, 0, 0, 255}
			}

			vector.DrawFilledRect(
				screen,
				float32(x+chunk.X*chunk.Game.ChunkWidth)*chunk.Game.CellSize+chunk.Game.SideBarLength,
				float32(y+chunk.Y*chunk.Game.ChunkHeight)*chunk.Game.CellSize,
				chunk.Game.CellSize,
				chunk.Game.CellSize,
				col,
				false,
			)
		}
	}
}

func (game *Game) Update() error {
	if err := game.ElementScrollBar.Update(); err != nil {
		return err
	}
	if cell, err := game.GetHoveredCell(); err == nil && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && game.SelectedElement != -1 {
		cell.Type = game.SelectedElement
	}
	return nil
}

func (game *Game) GetChunk(chunkX, chunkY int) (*Chunk, error) {
	if chunkX < 0 || chunkY < 0 || chunkX >= game.Width || chunkY >= game.Height {
		return nil, fmt.Errorf("there is no chunk at chunk position %v %v", chunkX, chunkY)
	}
	return &game.Chunks[game.CalculateChunkIndex(chunkX, chunkY)], nil
}

func (game *Game) GetCell(worldX, worldY int) (*Cell, error) {
	chunkX := worldX / game.ChunkWidth
	chunkY := worldY / game.ChunkHeight

	if chunk, err := game.GetChunk(chunkX, chunkY); err != nil {
		return nil, fmt.Errorf("there is no cell at world position %v %v", worldX, worldY)
	} else {
		cellX := worldX % game.ChunkWidth
		cellY := worldY % game.ChunkHeight
		if cell, err := chunk.GetCell(cellX, cellY); err != nil {
			return nil, fmt.Errorf("there is no cell at world position %v %v", worldX, worldY)
		} else {
			return cell, nil
		}
	}
}

func (game *Game) GetHoveredCell() (*Cell, error) {
	mx, my := ebiten.CursorPosition()
	x := float32(mx)
	y := float32(my)
	if x < game.SideBarLength {
		return nil, fmt.Errorf("%v %v is not on board", x, y)
	}
	x -= game.SideBarLength
	xIndex := int(x / game.CellSize)
	yIndex := int(y / game.CellSize)
	return game.GetCell(xIndex, yIndex)
}

func (chunk *Chunk) GetCell(cellX, cellY int) (*Cell, error) {
	if cellX < 0 || cellY < 0 || cellX >= chunk.Game.ChunkWidth || cellY >= chunk.Game.ChunkHeight {
		return nil, fmt.Errorf("there is no cell in chunk at local position %v %v", cellX, cellY)
	}
	return &chunk.Cells[chunk.Game.CalculateCellIndex(cellX, cellY)], nil
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
