package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type ScrollBarBox struct {
	Border, Inner color.Color
	BorderSize    float32
}

type ScrollBarItem struct {
	Box          *ScrollBarBox
	Background   color.Color
	InnerPadding float32
	TextColor    color.Color
	Text         string
	Clicked      func(item *ScrollBarItem, i int) error
	BeforeDraw   func(item *ScrollBarItem, i int)
}

type ScrollBar struct {
	X, Width        float32
	ElementHeight   float32
	Padding         float32
	Scroll          float32
	Items           []ScrollBarItem
	BackgroundColor color.Color
}

func (s *ScrollBar) GetHovered() int {
	x, y := ebiten.CursorPosition()
	return s.GetHoveredItem(float32(x), float32(y))
}

func NewScrollBar(x, width, elementHeight, padding float32, background color.Color, capacity int) ScrollBar {
	bar := ScrollBar{}

	bar.X = x
	bar.Width = width
	bar.ElementHeight = elementHeight
	bar.Padding = padding
	bar.BackgroundColor = background

	bar.Items = make([]ScrollBarItem, 0, capacity)
	bar.Scroll = 0
	return bar
}

func (s *ScrollBar) AddItem(item ScrollBarItem) {
	s.Items = append(s.Items, item)
}

func (scrollBar *ScrollBar) GetHoveredItem(x, y float32) int {
	if x < scrollBar.X || x > scrollBar.X+scrollBar.Width {
		return -1
	}

	index := int(y) / int(scrollBar.ElementHeight+scrollBar.Padding)

	if index < 0 || index >= len(scrollBar.Items) {
		return -1
	}

	return index
}

func (scrollBar *ScrollBar) ClickAt(x, y float32) error {
	y -= scrollBar.Scroll
	y += scrollBar.Padding

	index := scrollBar.GetHoveredItem(x, y)
	if index == -1 {
		return nil
	}

	item := scrollBar.Items[index]
	return item.Clicked(&item, index)
}

func (scrollBar *ScrollBar) Move(amt float32, x float32) {
	if x <= scrollBar.X || x >= scrollBar.X+scrollBar.Width {
		return
	}

	totalHeight := scrollBar.ElementHeight * float32(len(scrollBar.Items))
	totalHeight += scrollBar.Padding * float32(len(scrollBar.Items)-1)

	if totalHeight < float32(Dimensions.Height) {
		scrollBar.Scroll = 0
		return
	}

	scrollBar.Scroll += amt

	maxScroll := totalHeight - float32(Dimensions.Height)
	if scrollBar.Scroll < 0 {
		scrollBar.Scroll = 0
	} else if scrollBar.Scroll > maxScroll {
		scrollBar.Scroll = maxScroll
	}
}

func (scrollBar *ScrollBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		scrollBar.X,
		0,
		scrollBar.Width,
		float32(Dimensions.Height),
		scrollBar.BackgroundColor,
		true,
	)

	y := -scrollBar.Scroll + scrollBar.Padding

	for i := range scrollBar.Items {
		x := scrollBar.X + scrollBar.Padding
		ly := y
		item := scrollBar.Items[i]

		item.BeforeDraw(&item, i)

		vector.DrawFilledRect(
			screen,
			x, y,
			scrollBar.Width-scrollBar.Padding*2,
			scrollBar.ElementHeight,
			item.Background,
			true,
		)

		x += item.InnerPadding
		ly += item.InnerPadding

		if item.Box != nil {
			vector.DrawFilledRect(
				screen,
				x, ly,
				scrollBar.ElementHeight-item.InnerPadding*2, scrollBar.ElementHeight-item.InnerPadding*2,
				item.Box.Border,
				true,
			)

			vector.DrawFilledRect(
				screen,
				x+item.Box.BorderSize, ly+item.Box.BorderSize,
				scrollBar.ElementHeight-item.Box.BorderSize*2-item.InnerPadding*2,
				scrollBar.ElementHeight-item.Box.BorderSize*2-item.InnerPadding*2,
				item.Box.Inner,
				true,
			)

			x += scrollBar.ElementHeight - item.InnerPadding*2 + scrollBar.Padding
		}

		drawOptions := text.DrawOptions{}
		drawOptions.GeoM.Translate(float64(x), float64(ly))

		r, g, b, a := item.TextColor.RGBA()

		drawOptions.ColorScale.SetR(float32(r))
		drawOptions.ColorScale.SetG(float32(g))
		drawOptions.ColorScale.SetB(float32(b))
		drawOptions.ColorScale.SetA(float32(a))

		drawOptions.LayoutOptions.LineSpacing = 13 + 3

		text.Draw(
			screen,
			item.Text,
			text.NewGoXFace(basicfont.Face7x13),
			&drawOptions,
		)

		y += scrollBar.ElementHeight + scrollBar.Padding
	}
}

func (scrollBar *ScrollBar) Update() error {
	x, y := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if err := scrollBar.ClickAt(float32(x), float32(y)); err != nil {
			return fmt.Errorf("error while selecting element: %v", err)
		}
	}

	_, dy := ebiten.Wheel()
	scrollBar.Move(float32(dy), float32(x))

	return nil
}
