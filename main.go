package main

import (
	"log"

	"github.com/guilledipa/7colors/sevencolors"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := sevencolors.NewGame()
	ebiten.SetWindowSize(sevencolors.ScreenWidth*2, sevencolors.ScreenHeight*2)
	ebiten.SetWindowTitle("7 Colors Clone")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
