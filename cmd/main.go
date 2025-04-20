package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"go-falling-sand/game"
)

func main() {
	ebiten.SetWindowTitle("Falling Sand Game")
	ebiten.SetWindowSize(game.Dimensions.ScreenWidth, game.Dimensions.ScreenHeight)

	if g, err := game.NewGame(8, 8, 10, 10, 5, 200, "./data"); err != nil {
		log.Fatal(err)
	} else if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
