package game

import (
	"encoding/xml"
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"go-falling-sand/xml_handler"

	"github.com/hajimehoshi/ebiten/v2"
)

var Dimensions = struct {
	Width, Height, ScreenWidth, ScreenHeight int
}{400, 400, 600, 600}

// HexStringToColor converts a hex color string (e.g., "#1A2B3C") to a color.Color.
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

type ElementData struct {
	Color           color.Color
	Name            string
	ElementTypeName string
	ElementTypeID   int
	IsDefault       bool
}

type Game struct {
	elementIdCounter int
	DefaultElement   int
	ElementTypes     map[string]int
	ElementData      map[int]ElementData
	Width, Height    int
	Cells            []Cell
}

func (g *Game) DefineElement(elementTypeName string, color string, name string, isDefault bool) error {
	var col, colErr = StringToColor(color)
	if colErr != nil {
		return colErr
	}

	g.ElementTypes[elementTypeName] = g.elementIdCounter
	g.ElementData[g.elementIdCounter] = ElementData{
		Color:           col,
		Name:            name,
		ElementTypeName: elementTypeName,
		ElementTypeID:   g.elementIdCounter,
		IsDefault:       isDefault,
	}

	if isDefault {
		g.DefaultElement = g.elementIdCounter
	}

	g.elementIdCounter++

	return nil
}

type Cell struct {
	Type int
	X, Y int
	Game *Game
}

func (game *Game) Area() int {
	return game.Width * game.Height
}

func NewGame(width, height int, xmlData []byte) (*Game, error) {
	game := &Game{}

	game.elementIdCounter = 0
	game.Width = width
	game.Height = height
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

		if err := game.DefineElement(command.Name, col, name, command.IsDefault); err != nil {
			return nil, err
		}
	}

	game.Cells = make([]Cell, game.Area())
	for x := 0; x < game.Width; x++ {
		for y := 0; y < game.Height; y++ {
			i := game.CalculateIndex(x, y)
			game.Cells[i] = Cell{
				X: x, Y: y,
				Type: game.DefaultElement,
				Game: game,
			}
		}
	}

	return game, nil
}

func (g *Game) CalculateIndex(x, y int) int {
	return x + y*g.Width
}

func (game *Game) Layout(outsizeWidth, outsizeHeight int) (int, int) {
	return Dimensions.Width, Dimensions.Height
}

func (game *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.Black)

}

func (game *Game) Update() error {

	return nil
}
