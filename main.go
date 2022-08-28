//go:build 7colors
// +build 7colors

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowTitle  = "7 colors"
	WindowWidth  = 640
	WindowHeight = 480
)

func main() {
	ebiten.SetWindowTitle(WindowTitle)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	// ebiten.SetWindowResizable(true)

	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err = ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
