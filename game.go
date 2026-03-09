package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

var gameRand *rand.Rand

func init() {
	gameRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewGame() *Game {
	g := &Game{}
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			g.board[y][x] = gameRand.Intn(len(palette))
		}
	}

	// El Jugador 1 arranca abajo a la izquierda (0, boardSize-1)
	g.playerOwned[0] = []Point{{X: 0, Y: boardSize - 1}}
	g.playerColors[0] = g.board[boardSize-1][0]

	// El Jugador 2 arranca arriba a la derecha (boardSize-1, 0)
	g.playerOwned[1] = []Point{{X: boardSize - 1, Y: 0}}
	g.playerColors[1] = g.board[0][boardSize-1]

	// Nos fijamos que no empiecen con el mismo color
	for g.playerColors[0] == g.playerColors[1] {
		g.board[0][boardSize-1] = gameRand.Intn(len(palette))
		g.playerColors[1] = g.board[0][boardSize-1]
	}

	// Captura inicial para los dos
	g.capture(0, g.playerColors[0])
	g.capture(1, g.playerColors[1])

	return g
}

func (g *Game) capture(playerIdx int, newColor int) {
	g.playerColors[playerIdx] = newColor

	// Cambiamos el color de todas las fichas que ya tenemos
	for _, p := range g.playerOwned[playerIdx] {
		g.board[p.Y][p.X] = newColor
	}

	// Buscamos nuevas fichas para capturar
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

	// === Lógica de Encierro (Block/Cluster formation) ===
	for checkPlayer := 0; checkPlayer < 2; checkPlayer++ {
		otherPlayer := 1 - checkPlayer

		otherOwned := make(map[Point]bool)
		for _, p := range g.playerOwned[otherPlayer] {
			otherOwned[p] = true
		}

		reachable := make(map[Point]bool)
		queue := []Point{}

		for _, p := range g.playerOwned[checkPlayer] {
			reachable[p] = true
			queue = append(queue, p)
		}

		if len(queue) == 0 {
			continue
		}

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			neighbors := []Point{
				{X: curr.X + 1, Y: curr.Y},
				{X: curr.X - 1, Y: curr.Y},
				{X: curr.X, Y: curr.Y + 1},
				{X: curr.X, Y: curr.Y - 1},
			}

			for _, n := range neighbors {
				if n.X >= 0 && n.X < boardSize && n.Y >= 0 && n.Y < boardSize {
					if !reachable[n] && !otherOwned[n] {
						reachable[n] = true
						queue = append(queue, n)
					}
				}
			}
		}

		// Todo lo del tablero que NO sea alcanzable y que tampoco sea ya del otro, fue encerrado
		for y := 0; y < boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				p := Point{X: x, Y: y}
				if !reachable[p] && !otherOwned[p] {
					g.playerOwned[otherPlayer] = append(g.playerOwned[otherPlayer], p)
					g.board[y][x] = g.playerColors[otherPlayer]
					otherOwned[p] = true 
				}
			}
		}
	}

	if len(g.playerOwned[0])+len(g.playerOwned[1]) == boardSize*boardSize {
		g.gameOver = true
		if len(g.playerOwned[0]) > len(g.playerOwned[1]) {
			g.winner = 0
		} else if len(g.playerOwned[1]) > len(g.playerOwned[0]) {
			g.winner = 1
		} else {
			g.winner = -1 // Empate
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		// Botones de la barra del costado
		for i := 0; i < len(palette); i++ {
			bx := 30
			by := 80 + i*55
			bw := 45
			bh := 45
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
