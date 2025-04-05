package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var Dimensions = struct {
	Width, Height, ScreenWidth, ScreenHeight int
}{400, 400, 600, 600}

func StringToColor(s string) (color.Color, error) {
	if len(s) == 0 {
		return color.Black, errors.New("can't give an empty string")
	}

	if s[0] != '#' {
		return color.Black, errors.New("TODO: emplement this")
	}

	if col, ok := ColorMap[s]; ok {
		return col, nil
	}

	return color.Black, errors.New("there is no color '" + s + "'")
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
