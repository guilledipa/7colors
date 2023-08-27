package sevencolors

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	boardSize    = 8
)

// score is used to track each player scores.
type score struct {
	Player1 int
	Player2 int
}

// Game represents a game state.
type Game struct {
	board         *Board
	boardImage    *ebiten.Image
	rng           *rand.Rand
	CurrentPlayer int
	currentTurn   color.Color
	score         *score
}

// NewGame generates a new Game object.
func NewGame() *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &Game{
		board:       NewBoard(boardSize, rng),
		currentTurn: generateRandomColor(rng),
		rng:         rng,
		score:       &score{},
	}
	return g
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Get mouse coordinates and convert them to grid coordinates
		mouseX, mouseY := ebiten.CursorPosition()
		gridX := mouseX / tileSize
		gridY := mouseY / tileSize
		// Check if the selected grid coordinates are within bounds
		if gridX >= 0 && gridX < g.board.size && gridY >= 0 && gridY < g.board.size {
			// Check if the clicked tile can be conquered
			selectedColor := g.board.grid[gridY][gridX]
			if selectedColor != g.currentTurn {
				// Implement color conquering logic here
				g.conquerTiles(gridX, gridY, selectedColor)
				// Update the CurrentTurn to the next player's color
				g.currentTurn = generateRandomColor(g.rng) // Or implement a different logic
			}
		}
	}

	// Check for winning condition here and handle game over if necessary

	return nil
}

// Draw draws the current game to the given screen, centered.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(g.board.Size())
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)

	// Draw the scores
	player1Score := fmt.Sprintf("Player 1 Score: %d", g.score.Player1)
	player2Score := fmt.Sprintf("Player 2 Score: %d", g.score.Player2)
	drawTextWithShadow(screen, player1Score, 10, 20, 2, color.Black)
	drawTextWithShadow(screen, player2Score, 10, 45, 2, color.Black)
}

func (g *Game) conquerTiles(gridX, gridY int, targetColor color.Color) {
	if gridX < 0 || gridX >= g.board.size || gridY < 0 || gridY >= g.board.size {
		return
	}
	currentColor := g.board.grid[gridY][gridX]
	if currentColor != targetColor {
		return
	}
	// Change the color of the current tile
	g.board.grid[gridY][gridX] = g.currentTurn
	// Recursively conquer adjacent tiles
	g.conquerTiles(gridX-1, gridY, targetColor)
	g.conquerTiles(gridX+1, gridY, targetColor)
	g.conquerTiles(gridX, gridY-1, targetColor)
	g.conquerTiles(gridX, gridY+1, targetColor)
}
