package sevencolors

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize   = 40
	tileMargin = 4
)

// Board represents the game board.
type Board struct {
	size        int
	grid        [][]color.Color
	currentTurn color.Color
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) *Board {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	board := &Board{
		size:        size,
		grid:        generateRandomGrid(size, rng),
		currentTurn: generateRandomColor(rng),
	}
	return board
}

func generateRandomGrid(gridSize int, rng *rand.Rand) [][]color.Color {
	grid := make([][]color.Color, gridSize)
	for i := 0; i < gridSize; i++ {
		grid[i] = make([]color.Color, gridSize)
		for j := 0; j < gridSize; j++ {
			grid[i][j] = generateRandomColor(rng)
		}
	}
	return grid
}

func (b *Board) Draw(screen *ebiten.Image) {
	screen.Fill(frameColor)
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			op := &ebiten.DrawImageOptions{}
			x := j*tileSize + (j+1)*tileMargin
			y := i*tileSize + (i+1)*tileMargin
			tileColor := b.grid[i][j]
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorScale.ScaleWithColor(tileColor)
			tileImage := ebiten.NewImage(tileSize, tileSize)
			tileImage.Fill(tileColor)
			screen.DrawImage(tileImage, op)
		}
	}
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}
