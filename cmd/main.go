package main

import (
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"go-falling-sand/game"
)

func main() {
	ebiten.SetWindowTitle("Falling Sand Game")
	ebiten.SetWindowSize(game.Dimensions.ScreenWidth, game.Dimensions.ScreenHeight)

	xmlFile, err := os.Open("data.xml")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	xmlFile.Close()

	if g, err := game.NewGame(8, 8, 10, 10, 5, 200, byteValue); err != nil {
		log.Fatal(err)
	} else if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
