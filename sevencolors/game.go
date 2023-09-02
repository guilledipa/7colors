package sevencolors

import (
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

var (
	keys = []ebiten.Key{
		ebiten.KeyR, // Red
		ebiten.KeyG, // Green
		ebiten.KeyB, // Blue
		ebiten.KeyY, // Yellow
		ebiten.KeyM, // Magenta
		ebiten.KeyC, // Cyan
		ebiten.KeyW, // White
	}
)

// Game represents a game state.
type Game struct {
	board         *Board
	boardImage    *ebiten.Image
	rng           *rand.Rand
	CurrentPlayer int
	currentTurn   color.Color
}

// NewGame generates a new Game object.
func NewGame() *Game {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := &Game{
		board:       NewBoard(boardSize, rng),
		currentTurn: generateRandomColor(rng),
		rng:         rng,
	}
	return g
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	for _, key := range keys {
		if !inpututil.IsKeyJustPressed(key) {
			continue
		}
		// conquer colors
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
}
