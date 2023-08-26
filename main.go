package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gridSize = 8
	tileSize = 40
)

var (
	backgroundColor = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
	frameColor      = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
)

type Game struct {
	Grid        [][]color.Color
	CurrentTurn color.Color
}

func NewGame(rng *rand.Rand) *Game {
	game := &Game{
		Grid:        generateRandomGrid(rng),
		CurrentTurn: sevencolors.generateRandomColor(rng),
	}
	return game
}

func generateRandomGrid(rng *rand.Rand) [][]color.Color {
	grid := make([][]color.Color, gridSize)
	for i := 0; i < gridSize; i++ {
		grid[i] = make([]color.Color, gridSize)
		for j := 0; j < gridSize; j++ {
			grid[i][j] = sevencolors.generateRandomColor(rng)
		}
	}
	return grid
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Handle tile selection and color changing logic here
	}

	// Check for winning condition and update game state

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(frameColor)
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			op := &ebiten.DrawImageOptions{}
			x := j * tileSize
			y := i * tileSize
			tileColor := g.Grid[i][j]
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(tileColor)
			tileImage := ebiten.NewImage(tileSize, tileSize)
			tileImage.Fill(tileColor)
			screen.DrawImage(tileImage, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("7 Colors Clone")

	game := NewGame(rng)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
