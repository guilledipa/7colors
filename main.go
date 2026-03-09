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
	// Paleta Vintent-DOS moderna
	palette = []color.RGBA{
		{25, 65, 140, 255},   // Azul marino
		{185, 45, 45, 255},   // Carmesí
		{35, 125, 65, 255},   // Esmeralda
		{225, 185, 30, 255},  // Oro
		{145, 55, 135, 255},  // Púrpura
		{50, 150, 160, 255},  // Turquesa
		{120, 130, 140, 255}, // Acero
	}
	fontSource *text.GoTextFaceSource
	diamondImg *ebiten.Image
	bevelTL    *ebiten.Image
	bevelTR    *ebiten.Image
	bevelBL    *ebiten.Image
	bevelBR    *ebiten.Image
	gameRand   *rand.Rand
)

func initSprites() {
	if diamondImg != nil {
		return
	}
	w := tileSize * 2
	h := tileSize

	// 1. Base plana
	diamondImg = ebiten.NewImage(w, h)
	path := vector.Path{}
	path.MoveTo(float32(tileSize), 0) // Superior
	path.LineTo(float32(tileSize*2), float32(tileSize/2)) // Derecha
	path.LineTo(float32(tileSize), float32(tileSize)) // Inferior
	path.LineTo(0, float32(tileSize/2)) // Izquierda
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = 1, 1, 1, 1
	}
	emptyImg := ebiten.NewImage(1, 1)
	emptyImg.Fill(color.White)
	diamondImg.DrawTriangles(vertices, indices, emptyImg, &ebiten.DrawTrianglesOptions{})

	// Función helper para crear biseles blancos opacos
	createBevel := func(p1, p2, p3 [2]float32) *ebiten.Image {
		img := ebiten.NewImage(w, h)
		p := vector.Path{}
		p.MoveTo(p1[0], p1[1])
		p.LineTo(p2[0], p2[1])
		p.LineTo(p3[0], p3[1])
		p.Close()
		
		v, ind := p.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range v {
			v[i].ColorR, v[i].ColorG, v[i].ColorB, v[i].ColorA = 1, 1, 1, 1
		}
		img.DrawTriangles(v, ind, emptyImg, &ebiten.DrawTrianglesOptions{})
		return img
	}

	cTop := [2]float32{float32(tileSize), 0}
	cRight := [2]float32{float32(tileSize * 2), float32(tileSize / 2)}
	cBottom := [2]float32{float32(tileSize), float32(tileSize)}
	cLeft := [2]float32{0, float32(tileSize / 2)}
	cCenter := [2]float32{float32(tileSize), float32(tileSize / 2)}

	bevelTL = createBevel(cTop, cCenter, cLeft)
	bevelTR = createBevel(cTop, cRight, cCenter)
	bevelBL = createBevel(cLeft, cCenter, cBottom)
	bevelBR = createBevel(cCenter, cRight, cBottom)
}

func init() {
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
	drawOp       ebiten.DrawImageOptions // Reutilizada para performance
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
	// Revisamos a ambos jugadores para ver si el otro le cortó el paso
	for checkPlayer := 0; checkPlayer < 2; checkPlayer++ {
		otherPlayer := 1 - checkPlayer

		// 1. Armamos un mapa con las fichas que ya son de "otherPlayer", porque por ahí no podemos pasar
		otherOwned := make(map[Point]bool)
		for _, p := range g.playerOwned[otherPlayer] {
			otherOwned[p] = true
		}

		// 2. Buscamos todas las fichas a las que checkPlayer tiene acceso (Flood Fill)
		reachable := make(map[Point]bool)
		queue := []Point{}

		// Empezamos desde todas las fichas que checkPlayer ya tiene
		for _, p := range g.playerOwned[checkPlayer] {
			reachable[p] = true
			queue = append(queue, p)
		}

		// Si el jugador no tiene fichas (ej. test), asumimos que lo único que lo frena son las fichas del otro
		if len(queue) == 0 {
			// En un juego normal no pasa, pero en testing sí. Le damos una semilla "teórica" fuera del tablero?
			// Mejor no hacer nada, si no tiene fichas, todo es inalcanzable, o sea, todo se lo regalaría al otro.
			// PERO ojo, no queremos regalar todo el tablero vacío así no más.
			continue
		}

		// Expandimos
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			neighbors := []Point{
				{X: curr.X + 1, Y: curr.Y},
				{X: curr.X - 1, Y: curr.Y},
				{X: curr.X, Y: curr.Y + 1},
				{X: curr.X, Y: curr.Y - 1},
			}

			// NOTA: Para no encerrar "el vacio", solo podemos movernos si estamos dentro del board,
			// PERO la regla original asume que puedes conectar contra los bordes.
			// Mi lógica asume que todo lo que no es tocado por tu floodfill es encerrado.
			
			for _, n := range neighbors {
				if n.X >= 0 && n.X < boardSize && n.Y >= 0 && n.Y < boardSize {
					// Si no está visitado y el OTRO jugador no es dueño, podemos pasar
					if !reachable[n] && !otherOwned[n] {
						reachable[n] = true
						queue = append(queue, n)
					}
				}
			}
		}

		// 3. Todo lo del tablero que NO sea alcanzable por checkPlayer, 
		// y que tampoco sea ya de otherPlayer...
		// Significa que otherPlayer lo encerró. ¡Se lo regalamos a otherPlayer!
		for y := 0; y < boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				p := Point{X: x, Y: y}
				if !reachable[p] && !otherOwned[p] {
					g.playerOwned[otherPlayer] = append(g.playerOwned[otherPlayer], p)
					g.board[y][x] = g.playerColors[otherPlayer]
					otherOwned[p] = true // Para no agregarlo de nuevo si hay superposición
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
	initSprites()
	// Fondo gris oscuro texturado para feeling DOS / Mármol
	screen.Fill(color.RGBA{20, 25, 30, 255})

	fs := getFontSource()

	// Título sombra + texto
	opShadow := &text.DrawOptions{}
	opShadow.GeoM.Translate(22, 22)
	opShadow.ColorScale.Scale(0, 0, 0, 0.8)
	text.Draw(screen, "7 COLORS", &text.GoTextFace{Source: fs, Size: 32}, opShadow)

	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	text.Draw(screen, "7 COLORS", &text.GoTextFace{Source: fs, Size: 32}, op)

	// Barra lateral para elegir color
	for i, clr := range palette {
		bx := float32(30)
		by := float32(80 + i*55)
		bw := float32(45)
		bh := float32(45)

		// Es jugador activo?
		isPlayerChoice := i == g.playerColors[0] || i == g.playerColors[1]

		// Sombra del botón
		vector.DrawFilledRect(screen, bx+3, by+3, bw, bh, color.RGBA{0, 0, 0, 150}, false)

		if isPlayerChoice {
			// Dim if it's taken
			dimClr := color.RGBA{uint8(float32(clr.R) * 0.4), uint8(float32(clr.G) * 0.4), uint8(float32(clr.B) * 0.4), 255}
			vector.DrawFilledRect(screen, bx, by, bw, bh, dimClr, false)
			vector.StrokeRect(screen, bx, by, bw, bh, 2, color.RGBA{100, 100, 100, 255}, false)

			var playerClr color.Color
			var playerText string
			if i == g.playerColors[0] {
				playerClr = color.White
				playerText = "P1"
			} else {
				playerClr = color.RGBA{170, 170, 170, 255}
				playerText = "P2"
			}

			// Marca de que está tomado
			vector.StrokeLine(screen, bx+5, by+5, bx+bw-5, by+bh-5, 2, color.RGBA{0, 0, 0, 200}, false)
			vector.StrokeLine(screen, bx+bw-5, by+5, bx+5, by+bh-5, 2, color.RGBA{0, 0, 0, 200}, false)

			// Outline indicator for text
			opShadowL := &text.DrawOptions{}
			opShadowL.GeoM.Translate(float64(bx+bw/4)+1, float64(by+bh/4)+1)
			opShadowL.ColorScale.Scale(0, 0, 0, 1)
			text.Draw(screen, playerText, &text.GoTextFace{Source: fs, Size: 16}, opShadowL)

			// Text
			opLabel := &text.DrawOptions{}
			opLabel.GeoM.Translate(float64(bx+bw/4), float64(by+bh/4))
			r, g, b, _ := playerClr.RGBA()
			opLabel.ColorScale.Scale(float32(r)/65535, float32(g)/65535, float32(b)/65535, 1)
			text.Draw(screen, playerText, &text.GoTextFace{Source: fs, Size: 16}, opLabel)

		} else {
			// Clickable button
			vector.DrawFilledRect(screen, bx, by, bw, bh, clr, false)
			
			// Highlight top-left border
			vector.StrokeLine(screen, bx, by, bx+bw, by, 2, color.RGBA{255, 255, 255, 120}, false)
			vector.StrokeLine(screen, bx, by, bx, by+bh, 2, color.RGBA{255, 255, 255, 120}, false)
			// Dark bottom-right border
			vector.StrokeLine(screen, bx+bw, by, bx+bw, by+bh, 2, color.RGBA{0, 0, 0, 120}, false)
			vector.StrokeLine(screen, bx, by+bh, bx+bw, by+bh, 2, color.RGBA{0, 0, 0, 120}, false)
		}
	}

	// Dibujamos los puntajes
	totalTiles := float64(boardSize * boardSize)
	p1Score := float64(len(g.playerOwned[0])) / totalTiles * 100
	p2Score := float64(len(g.playerOwned[1])) / totalTiles * 100

	drawShadowText := func(scr *ebiten.Image, txt string, size float64, fX, fY float64, clr color.Color) {
		// Sombra
		opS := &text.DrawOptions{}
		opS.GeoM.Translate(fX+2, fY+2)
		opS.ColorScale.Scale(0, 0, 0, 0.7)
		text.Draw(scr, txt, &text.GoTextFace{Source: fs, Size: size}, opS)
		// Texto
		opT := &text.DrawOptions{}
		opT.GeoM.Translate(fX, fY)
		r, gg, b, _ := clr.RGBA()
		opT.ColorScale.Scale(float32(r)/65535, float32(gg)/65535, float32(b)/65535, 1)
		text.Draw(scr, txt, &text.GoTextFace{Source: fs, Size: size}, opT)
	}

	drawShadowText(screen, fmt.Sprintf("P1 AREA: %.1f%%", p1Score), 24, screenWidth-220, 40, color.White)
	drawShadowText(screen, fmt.Sprintf("P2 AREA: %.1f%%", p2Score), 24, screenWidth-220, 80, color.RGBA{180, 180, 180, 255})

	// Turno de quién es
	turnText := ">>> Turno del Jugador 1 <<<"
	var turnColor color.Color = color.White
	if g.turn == 1 {
		turnText = ">>> Turno del Jugador 2 <<<"
		turnColor = color.RGBA{180, 180, 180, 255}
	}
	if g.gameOver {
		if g.winner == -1 {
			turnText = "!! EMPATE !! Clic para reiniciar"
			turnColor = color.RGBA{255, 255, 0, 255}
		} else {
			turnText = fmt.Sprintf("!! GANÓ EL JUGADOR %d !! Clic para reiniciar", g.winner+1)
			if g.winner == 0 {
				turnColor = color.RGBA{100, 255, 100, 255}
			} else {
				turnColor = color.RGBA{255, 100, 100, 255}
			}
		}
	}
	
	// Centrado chapuzero
	textW := float64(len(turnText)) * 11.5 // Aproximación
	drawShadowText(screen, turnText, 24, screenWidth/2-textW/2, screenHeight-50, turnColor)

	// Dibujamos el tablero propiamente dicho
	offsetX := float32(screenWidth) / 2 + 50
	offsetY := float32(screenHeight) / 2 - 30

	// PASADA 1: Dibujamos las bases planas de todos los diamantes
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			cx := float64(offsetX + float32(x-y)*float32(tileSize))
			cy := float64(offsetY + float32(x+y)*float32(tileSize)/2 - float32(boardSize)*float32(tileSize)/4)

			drawOp := &ebiten.DrawImageOptions{}
			drawOp.GeoM.Translate(cx-float64(tileSize), cy-float64(tileSize/2))

			clr := palette[g.board[y][x]]
			drawOp.ColorScale.Scale(float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255)
			screen.DrawImage(diamondImg, drawOp)
		}
	}

	// PASADA 2: Dibujamos los biseles 3D (Highlights/Shadows) solo en los bordes exteriores
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			clrIdx := g.board[y][x]
			baseClr := palette[clrIdx]

			// Calculamos colores claros y oscuros sólidos para los biseles
			r, gg, b := float32(baseClr.R)/255.0, float32(baseClr.G)/255.0, float32(baseClr.B)/255.0

			lR, lG, lB := r*1.4, gg*1.4, b*1.4
			if lR > 1.0 { lR = 1.0 }
			if lG > 1.0 { lG = 1.0 }
			if lB > 1.0 { lB = 1.0 }

			dR, dG, dB := r*0.6, gg*0.6, b*0.6

			cx := float64(offsetX + float32(x-y)*float32(tileSize))
			cy := float64(offsetY + float32(x+y)*float32(tileSize)/2 - float32(boardSize)*float32(tileSize)/4)

			drawOp := &ebiten.DrawImageOptions{}
			drawOp.GeoM.Translate(cx-float64(tileSize), cy-float64(tileSize/2))

			// Vecino Bottom-Right (X+1, Y). Si es diferente o es borde del tablero, dibujamos la sombra BR
			if x+1 >= boardSize || g.board[y][x+1] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(dR, dG, dB, 1.0)
				screen.DrawImage(bevelBR, op)
			}
			// Vecino Bottom-Left (X, Y+1). Si es diferente o es borde del tablero, dibujamos la sombra BL
			if y+1 >= boardSize || g.board[y+1][x] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(dR, dG, dB, 1.0)
				screen.DrawImage(bevelBL, op)
			}
			// Vecino Top-Left (X-1, Y). Si es diferente o es borde del tablero, dibujamos la luz TL
			if x-1 < 0 || g.board[y][x-1] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(lR, lG, lB, 1.0)
				screen.DrawImage(bevelTL, op)
			}
			// Vecino Top-Right (X, Y-1). Si es diferente o es borde del tablero, dibujamos la luz TR
			if y-1 < 0 || g.board[y-1][x] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(lR, lG, lB, 1.0)
				screen.DrawImage(bevelTR, op)
			}
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
