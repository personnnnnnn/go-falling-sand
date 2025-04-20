package game

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go-falling-sand/util"
	"go-falling-sand/xml_handler"

	"github.com/hajimehoshi/ebiten/v2"
)

var Dimensions = struct {
	Width, Height, ScreenWidth, ScreenHeight int
}{600, 400, 900, 600}

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
	OtherKinds      []ElementKind
	Bouyancy        float32
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
	ChunkOrder              []int
	CellSize                float32
	ElementScrollBar        ScrollBar
	UpdateCycle             bool
}

func (g *Game) TotalWidth() int {
	return g.Width * g.ChunkWidth
}

func (g *Game) TotalHeight() int {
	return g.Height * g.ChunkHeight
}

func (g *Game) DefineElement(
	definition *xmlhandler.XMLElementDefinition,
	elementTypeName string,
	colorString string,
	name string,
	role string,
	selectable bool,
	bouyancy float32,
) error {
	index := g.elementIdCounter
	g.elementIdCounter++

	var col, colErr = StringToColor(colorString)
	if colErr != nil {
		return colErr
	}

	if role != ROLE_AIR && role != ROLE_WALL && role != ROLE_NONE {
		role = ROLE_NONE
	}

	var kind ElementKind = nil

	if definition.MovableSolid != nil {
		kind = &MovableSolid{}
	} else if definition.Liquid != nil {
		kind = &Liquid{}
	} else if definition.Gas != nil {
		kind = &Gas{definition.Gas.Weight}
	} else if definition.Dust != nil {
		kind = &Dust{definition.Dust.Weight}
	} else if definition.ImmovableSolid != nil {
		kind = &ImmovableSolid{}
	} else {
		kind = &DefaultKind{}
	}

	g.ElementTypes[elementTypeName] = index
	g.ElementData[index] = ElementData{
		Color:           col,
		Name:            name,
		ElementTypeName: elementTypeName,
		ElementTypeID:   index,
		Role:            role,
		Bouyancy:        bouyancy,
		Kind:            kind,
		OtherKinds:      make([]ElementKind, 0, 2),
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

func (g *Game) DefineTransformations(definiton *xmlhandler.XMLElementDefinition) error {
	index := g.ElementTypes[definiton.Name]
	if definiton.Reactions != nil {
		for _, reaction := range definiton.Reactions.Reactions {
			kind := &Reaction{
				Actions:    make([]Action, 0, 2),
				Conditions: make([]Condition, 0, 2),
			}
			for _, v := range reaction.Steps {
				switch v.XMLName.Local {
				case "turn-into":
					{
						if id, ok := g.ElementTypes[v.Value]; !ok {
							return fmt.Errorf("there is no element named '%v'", v.Value)
						} else {
							kind.Actions = append(kind.Actions, &TurnInto{id})
						}
					}
				case "emit":
					{
						if id, ok := g.ElementTypes[v.Value]; !ok {
							return fmt.Errorf("there is no element named '%v'", v.Value)
						} else {
							kind.Actions = append(kind.Actions, &Emit{id})
						}
					}
				case "chance":
					{
						if chance, err := strconv.ParseFloat(v.Value, 32); err != nil {
							return err
						} else {
							kind.Conditions = append(kind.Conditions, &Chance{float32(chance)})
						}
					}
				case "touching":
					{
						if id, ok := g.ElementTypes[v.Value]; !ok {
							return fmt.Errorf("there is no element named '%v'", v.Value)
						} else {
							kind.Conditions = append(kind.Conditions, &Touching{id})
						}
					}
				case "directly-touching":
					{
						if id, ok := g.ElementTypes[v.Value]; !ok {
							return fmt.Errorf("there is no element named '%v'", v.Value)
						} else {
							kind.Conditions = append(kind.Conditions, &DirectlyTouching{id})
						}
					}
				}
			}
			elementData := g.ElementData[index]
			elementData.OtherKinds = append(elementData.OtherKinds, kind)
			g.ElementData[index] = elementData
		}
	}
	return nil
}

func (game *Game) ChunkArea() int {
	return game.ChunkWidth * game.ChunkHeight
}

func (game *Game) WorldArea() int {
	return game.Width * game.Height
}

func (game *Game) HandleCommand(command *xmlhandler.XMLElementDefinition) error {
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

	material := command.Material
	if material == nil {
		material = &xmlhandler.XMLMaterialData{}
	}

	if err := game.DefineElement(command, command.Name, col, name, command.Role, display.Selectable, material.Density); err != nil {
		return err
	}
	return nil
}

func (game *Game) HandleCommandReaction(command *xmlhandler.XMLElementDefinition) error {
	err := game.DefineTransformations(command)
	if err != nil {
		return err
	}
	return nil
}

func NewGame(width, height int, chunkWidth, chunkHeight int, cellSize float32, sideBarLength float32, dataFolder string) (*Game, error) {
	game := &Game{}

	game.elementIdCounter = 0

	game.SelectedElement = -1 // No item selected

	game.UpdateCycle = false

	game.Width = width
	game.Height = height

	game.ChunkWidth = chunkWidth
	game.ChunkHeight = chunkHeight

	game.CellSize = cellSize
	game.SideBarLength = sideBarLength

	game.ElementData = map[int]ElementData{}
	game.ElementTypes = map[string]int{}

	game.ElementScrollBar = NewScrollBar(
		0,
		sideBarLength,
		30,
		10,
		color.RGBA{100, 100, 100, 255},
		20,
	)

	matches := make([]string, 0, 20)
	err := filepath.WalkDir(dataFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".xml") {
			matches = append(matches, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error while getting xml files: %v", err)
	}

	results := make([]xmlhandler.XMLElementDefinition, 0, len(matches))

	for _, file := range matches {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file '%s': %v", file, err)
		}

		var elem xmlhandler.XMLElementDefinition
		if err := xml.Unmarshal(data, &elem); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file '%s': %v", file, err)
		}

		results = append(results, elem)
	}

	for _, result := range results {
		game.HandleCommand(&result)
	}

	for _, result := range results {
		game.HandleCommandReaction(&result)
	}

	game.Chunks = make([]Chunk, game.WorldArea())
	game.ChunkOrder = make([]int, game.WorldArea())
	for x := range game.Width {
		for y := range game.Height {
			i := game.CalculateChunkIndex(x, y)
			game.ChunkOrder[i] = i
			chunk := NewChunk(game, x, y)
			game.Chunks[i] = chunk
		}
	}

	return game, nil
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

			chunk := game.Chunks[i]
			chunk.Draw(screen)
		}
	}

	game.ElementScrollBar.Draw(screen)
}

func (game *Game) UpdateChunks() error {
	util.Shuffle(game.ChunkOrder)
	for i := range game.ChunkOrder {
		i = game.ChunkOrder[i]
		chunk := game.Chunks[i]
		if err := chunk.Update(); err != nil {
			return err
		}
	}
	game.UpdateCycle = !game.UpdateCycle
	return nil
}

func (game *Game) Update() error {
	if err := game.ElementScrollBar.Update(); err != nil {
		return err
	}
	if err := game.UpdateChunks(); err != nil {
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
