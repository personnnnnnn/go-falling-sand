package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"go-falling-sand/game"
)

type XMLElementList struct {
	XMLName  xml.Name               `xml:"elements"`
	Elements []XMLElementDefinition `xml:"element"`
}

type XMLElementDefinition struct {
	XMLName xml.Name    `xml:"element"`
	Name    string      `xml:"name,attr"`
	Display *XMLDisplay `xml:"display"`
}

type XMLDisplay struct {
	XMLName xml.Name `xml:"display"`
	Name    string   `xml:"name"`
	Color   string   `xml:"color"`
}

func main() {
	ebiten.SetWindowTitle("Pong in Ebiten")
	ebiten.SetWindowSize(game.Dimensions.ScreenWidth, game.Dimensions.ScreenHeight)

	xmlFile, err := os.Open("data.xml")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := io.ReadAll(xmlFile)

	var commands struct{}
	xml.Unmarshal(byteValue, &commands)

	xmlFile.Close()

	g := &game.Game{}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
