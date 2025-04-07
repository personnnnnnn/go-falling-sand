package game

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
)

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
