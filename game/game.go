package game

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

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
}

type Game struct {
	elementIdCounter int
	ElementTypes     map[string]int
	ElementData      map[int]ElementData
	Width, Height    int
	Cells            []Cell
}

func (g *Game) DefineElement(elementTypeName string, color string, name string) error {
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
	}

	g.elementIdCounter++

	return nil
}

type Cell struct {
	Type int
	X, Y int
	Game *Game
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
