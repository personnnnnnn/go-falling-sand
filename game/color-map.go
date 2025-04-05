package game

import "image/color"

var ColorMap = map[string]color.Color{
	"black":    color.RGBA{0, 0, 0, 255},
	"white":    color.RGBA{255, 255, 255, 255},
	"gray":     color.RGBA{100, 100, 100, 255},
	"grey":     color.RGBA{100, 100, 100, 255},
	"red":      color.RGBA{255, 0, 0, 255},
	"green":    color.RGBA{0, 255, 0, 255},
	"deepblue": color.RGBA{0, 0, 255, 255},
	"blue":     color.RGBA{0, 196, 255, 255},
	"cyan":     color.RGBA{0, 255, 255, 255},
	"yellow":   color.RGBA{255, 255, 0, 255},
	"sun":      color.RGBA{255, 194, 0, 255},
	"orange":   color.RGBA{255, 109, 0, 255},
	"rose":     color.RGBA{255, 0, 127, 255},
	"lavender": color.RGBA{255, 0, 220, 255},
	"magenta":  color.RGBA{255, 0, 255, 255},
	"sakura":   color.RGBA{255, 171, 255, 255},
	"hotpink":  color.RGBA{255, 68, 178, 255},
	"pink":     color.RGBA{255, 141, 212, 255},
	"purpur":   color.RGBA{157, 110, 253, 255},
	"purple":   color.RGBA{106, 69, 178, 255},
	"lime":     color.RGBA{106, 255, 178, 255},
	"brown":    color.RGBA{106, 69, 0, 255},
}
