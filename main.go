package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	screenWidth  = 640
	screenHeight = 480
	boardSize    = 13 // 13x13 grid of diamonds
	tileSize     = 20 // Adjusted size
)

var (
	// DOS-like palette for 7 Colors
	palette = []color.RGBA{
		{0, 0, 170, 255},     // Dark Blue
		{170, 0, 0, 255},     // Dark Red
		{0, 170, 0, 255},     // Dark Green
		{255, 255, 0, 255},   // Yellow
		{170, 0, 170, 255},   // Magenta
		{0, 170, 170, 255},   // Cyan
		{170, 170, 170, 255}, // Light Grey
	}
	fontSource *text.GoTextFaceSource
)

func init() {
	// Only initialize font if we're not in a test (to avoid issues in CI environments)
	// Actually, let's keep it but handle the error better.
}

func getFontSource() *text.GoTextFaceSource {
	if fontSource == nil {
		s, err := text.NewGoTextFaceSourceFromReader(bytes.NewReader(goregular.TTF))
		if err != nil {
			log.Fatal(err)
		}
		fontSource = s
	}
	return fontSource
}

type Point struct {
	X, Y int
}

type Game struct {
	board        [boardSize][boardSize]int
	playerColors [2]int
	playerOwned  [2][]Point
	turn         int
	gameOver     bool
	winner       int
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{}
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			g.board[y][x] = rand.Intn(len(palette))
		}
	}

	// Player 1 starts at bottom-left (0, boardSize-1)
	g.playerOwned[0] = []Point{{X: 0, Y: boardSize - 1}}
	g.playerColors[0] = g.board[boardSize-1][0]

	// Player 2 starts at top-right (boardSize-1, 0)
	g.playerOwned[1] = []Point{{X: boardSize - 1, Y: 0}}
	g.playerColors[1] = g.board[0][boardSize-1]

	// Ensure starting colors are different
	for g.playerColors[0] == g.playerColors[1] {
		g.board[0][boardSize-1] = rand.Intn(len(palette))
		g.playerColors[1] = g.board[0][boardSize-1]
	}

	// Initial capture for both players
	g.capture(0, g.playerColors[0])
	g.capture(1, g.playerColors[1])

	return g
}

func (g *Game) capture(playerIdx int, newColor int) {
	g.playerColors[playerIdx] = newColor

	// Update all currently owned tiles to the new color
	for _, p := range g.playerOwned[playerIdx] {
		g.board[p.Y][p.X] = newColor
	}

	// Find new tiles to capture
	for {
		added := false
		ownedMap := make(map[Point]bool)
		for _, p := range g.playerOwned[playerIdx] {
			ownedMap[p] = true
		}

		otherPlayerOwnedMap := make(map[Point]bool)
		for _, p := range g.playerOwned[1-playerIdx] {
			otherPlayerOwnedMap[p] = true
		}

		var newTiles []Point
		for _, p := range g.playerOwned[playerIdx] {
			neighbors := []Point{
				{X: p.X + 1, Y: p.Y},
				{X: p.X - 1, Y: p.Y},
				{X: p.X, Y: p.Y + 1},
				{X: p.X, Y: p.Y - 1},
			}

			for _, n := range neighbors {
				if n.X >= 0 && n.X < boardSize && n.Y >= 0 && n.Y < boardSize {
					if !ownedMap[n] && !otherPlayerOwnedMap[n] && g.board[n.Y][n.X] == newColor {
						newTiles = append(newTiles, n)
						ownedMap[n] = true
						added = true
					}
				}
			}
		}
		if !added {
			break
		}
		g.playerOwned[playerIdx] = append(g.playerOwned[playerIdx], newTiles...)
	}

	if len(g.playerOwned[0])+len(g.playerOwned[1]) == boardSize*boardSize {
		g.gameOver = true
		if len(g.playerOwned[0]) > len(g.playerOwned[1]) {
			g.winner = 0
		} else if len(g.playerOwned[1]) > len(g.playerOwned[0]) {
			g.winner = 1
		} else {
			g.winner = -1 // Tie
		}
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			*g = *NewGame()
		}
		return nil
	}

	// Mouse click handling for colors
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		// Sidebar buttons
		for i := 0; i < len(palette); i++ {
			bx := 20
			by := 50 + i*60
			bw := 50
			bh := 50
			if mx >= bx && mx <= bx+bw && my >= by && my <= by+bh {
				if i != g.playerColors[0] && i != g.playerColors[1] {
					g.capture(g.turn, i)
					g.turn = 1 - g.turn
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Sidebar for color choice
	for i, clr := range palette {
		bx := float32(20)
		by := float32(50 + i*60)
		bw := float32(50)
		bh := float32(50)

		// Button
		vector.DrawFilledRect(screen, bx, by, bw, bh, clr, false)

		// If color is selected by a player, mark it
		if i == g.playerColors[0] || i == g.playerColors[1] {
			var playerClr color.Color = color.White
			if i == g.playerColors[0] {
				playerClr = color.RGBA{255, 255, 255, 255}
			} else {
				playerClr = color.RGBA{100, 100, 100, 255}
			}
			vector.StrokeRect(screen, bx-2, by-2, bw+4, bh+4, 2, playerClr, false)
		}
	}

	// Draw scores
	totalTiles := float64(boardSize * boardSize)
	p1Score := float64(len(g.playerOwned[0])) / totalTiles * 100
	p2Score := float64(len(g.playerOwned[1])) / totalTiles * 100

	fs := getFontSource()
	op := &text.DrawOptions{}
	op.GeoM.Translate(100, 20)
	text.Draw(screen, fmt.Sprintf("Player 1: %.1f%%", p1Score), &text.GoTextFace{Source: fs, Size: 20}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth-200, 20)
	text.Draw(screen, fmt.Sprintf("Player 2: %.1f%%", p2Score), &text.GoTextFace{Source: fs, Size: 20}, op)

	// Current Turn
	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2-100, screenHeight-40)
	turnText := "Player 1's Turn"
	if g.turn == 1 {
		turnText = "Player 2's Turn"
	}
	if g.gameOver {
		if g.winner == -1 {
			turnText = "Tie Game! Click to restart."
		} else {
			turnText = fmt.Sprintf("Player %d Wins! Click to restart.", g.winner+1)
		}
	}
	text.Draw(screen, turnText, &text.GoTextFace{Source: fs, Size: 24}, op)

	// Draw board
	offsetX := float32(screenWidth) / 2 + 50
	offsetY := float32(screenHeight) / 2

	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			cx := offsetX + float32(x-y)*float32(tileSize)
			cy := offsetY + float32(x+y)*float32(tileSize)/2 - float32(boardSize)*float32(tileSize)/4

			g.drawDiamond(screen, cx, cy, palette[g.board[y][x]])
		}
	}
}

func (g *Game) drawDiamond(screen *ebiten.Image, cx, cy float32, clr color.Color) {
	w := float32(tileSize)
	h := float32(tileSize) / 2

	path := vector.Path{}
	path.MoveTo(cx, cy-h)
	path.LineTo(cx+w, cy)
	path.LineTo(cx, cy+h)
	path.LineTo(cx-w, cy)
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		r, g, b, a := clr.RGBA()
		vertices[i].SrcX = 0
		vertices[i].SrcY = 0
		vertices[i].ColorR = float32(r) / 0xffff
		vertices[i].ColorG = float32(g) / 0xffff
		vertices[i].ColorB = float32(b) / 0xffff
		vertices[i].ColorA = float32(a) / 0xffff
	}

	emptyImage := ebiten.NewImage(1, 1)
	emptyImage.Fill(color.White)
	screen.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("7 Colors - DOS Clone")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
