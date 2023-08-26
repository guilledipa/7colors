package main

import (
	"log"

	"github.com/guilledipa/7colors/sevencolors"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := sevencolors.NewGame()
	ebiten.SetWindowSize(sevencolors.ScreenWidth, sevencolors.ScreenHeight)
	ebiten.SetWindowTitle("7 Colors Clone")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
