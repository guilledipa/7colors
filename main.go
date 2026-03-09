package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	boardSize    = 13 // Tablero de 13x13 diamantes
	tileSize     = 20 // Tamaño ajustado
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("7 COLORS - CLON DE DOS")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
