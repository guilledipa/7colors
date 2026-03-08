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
	boardSize    = 13 // Tablero de 13x13 diamantes
	tileSize     = 20 // Tamaño ajustado
)

var (
	// Paleta tipo DOS para 7 Colors
	palette = []color.RGBA{
		{0, 0, 170, 255},     // Azul oscuro
		{170, 0, 0, 255},     // Rojo oscuro
		{0, 170, 0, 255},     // Verde oscuro
		{255, 255, 0, 255},   // Amarillo
		{170, 0, 170, 255},   // Magenta
		{0, 170, 170, 255},   // Cian
		{170, 170, 170, 255}, // Gris claro
	}
	fontSource *text.GoTextFaceSource
	diamondImg *ebiten.Image
	gameRand   *rand.Rand
)

func init() {
	// Armamos el sprite del diamante
	w := tileSize * 2
	h := tileSize
	diamondImg = ebiten.NewImage(w, h)

	path := vector.Path{}
	path.MoveTo(float32(tileSize), 0)
	path.LineTo(float32(tileSize*2), float32(tileSize/2))
	path.LineTo(float32(tileSize), float32(tileSize))
	path.LineTo(0, float32(tileSize/2))
	path.Close()

	// Lo rellenamos con blanco (así después le mandamos el color que queramos)
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].ColorR = 1
		vertices[i].ColorG = 1
		vertices[i].ColorB = 1
		vertices[i].ColorA = 1
	}
	emptyImg := ebiten.NewImage(1, 1)
	emptyImg.Fill(color.White)
	diamondImg.DrawTriangles(vertices, indices, emptyImg, &ebiten.DrawTrianglesOptions{})

	// Le mandamos un contorno para que quede más prolijo
	vector.StrokeLine(diamondImg, float32(tileSize), 0, float32(tileSize*2), float32(tileSize/2), 1, color.RGBA{0, 0, 0, 150}, false)
	vector.StrokeLine(diamondImg, float32(tileSize*2), float32(tileSize/2), float32(tileSize), float32(tileSize), 1, color.RGBA{0, 0, 0, 150}, false)
	vector.StrokeLine(diamondImg, float32(tileSize), float32(tileSize), 0, float32(tileSize/2), 1, color.RGBA{0, 0, 0, 150}, false)
	vector.StrokeLine(diamondImg, 0, float32(tileSize/2), float32(tileSize), 0, 1, color.RGBA{0, 0, 0, 150}, false)

	gameRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getFontSource() *text.GoTextFaceSource {
	if fontSource == nil {
		s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
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

	// Manejo de clics para elegir los colores
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

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	fs := getFontSource()

	// Título
	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	text.Draw(screen, "7 COLORS", &text.GoTextFace{Source: fs, Size: 32}, op)

	// Barra lateral para elegir color
	for i, clr := range palette {
		bx := float32(30)
		by := float32(80 + i*55)
		bw := float32(45)
		bh := float32(45)

		// Botón
		vector.DrawFilledRect(screen, bx, by, bw, bh, clr, false)
		// Borde para el botón
		vector.StrokeRect(screen, bx, by, bw, bh, 1, color.RGBA{255, 255, 255, 100}, false)

		// Si algún jugador ya tiene este color, lo marcamos
		if i == g.playerColors[0] || i == g.playerColors[1] {
			var playerClr color.Color = color.White
			var playerText string
			if i == g.playerColors[0] {
				playerClr = color.RGBA{255, 255, 255, 255}
				playerText = "P1"
			} else {
				playerClr = color.RGBA{170, 170, 170, 255}
				playerText = "P2"
			}
			vector.StrokeRect(screen, bx-3, by-3, bw+6, bh+6, 2, playerClr, false)

			// Etiqueta para saber de quién es
			opLabel := &text.DrawOptions{}
			opLabel.GeoM.Translate(float64(bx+bw/4), float64(by+bh/4))
			text.Draw(screen, playerText, &text.GoTextFace{Source: fs, Size: 16}, opLabel)
		}
	}

	// Dibujamos los puntajes
	totalTiles := float64(boardSize * boardSize)
	p1Score := float64(len(g.playerOwned[0])) / totalTiles * 100
	p2Score := float64(len(g.playerOwned[1])) / totalTiles * 100

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth-220, 40)
	text.Draw(screen, fmt.Sprintf("P1: %.1f%%", p1Score), &text.GoTextFace{Source: fs, Size: 24}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth-220, 80)
	text.Draw(screen, fmt.Sprintf("P2: %.1f%%", p2Score), &text.GoTextFace{Source: fs, Size: 24}, op)

	// Turno de quién es
	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2-120, screenHeight-50)
	turnText := ">>> Turno del Jugador 1 <<<"
	if g.turn == 1 {
		turnText = ">>> Turno del Jugador 2 <<<"
	}
	if g.gameOver {
		if g.winner == -1 {
			turnText = "!! EMPATE !! Clic para reiniciar"
		} else {
			turnText = fmt.Sprintf("!! GANÓ EL JUGADOR %d !! Clic para reiniciar", g.winner+1)
		}
	}
	text.Draw(screen, turnText, &text.GoTextFace{Source: fs, Size: 24}, op)

	// Dibujamos el tablero propiamente dicho
	offsetX := float32(screenWidth) / 2 + 60
	offsetY := float32(screenHeight) / 2 - 20

	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			cx := offsetX + float32(x-y)*float32(tileSize)
			cy := offsetY + float32(x+y)*float32(tileSize)/2 - float32(boardSize)*float32(tileSize)/4

			// Dibujamos el diamante con el sprite que ya armamos
			drawOp := &ebiten.DrawImageOptions{}
			drawOp.GeoM.Translate(float64(cx-float32(tileSize)), float64(cy-float32(tileSize/2)))

			clr := palette[g.board[y][x]]
			drawOp.ColorScale.Scale(float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255)

			screen.DrawImage(diamondImg, drawOp)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("7 COLORS - CLON DE DOS")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
