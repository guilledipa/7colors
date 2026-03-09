package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
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

	// Sprites
	diamondImg *ebiten.Image
	bevelTL    *ebiten.Image
	bevelTR    *ebiten.Image
	bevelBL    *ebiten.Image
	bevelBR    *ebiten.Image
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
	path.MoveTo(float32(tileSize), 0)
	path.LineTo(float32(tileSize*2), float32(tileSize/2))
	path.LineTo(float32(tileSize), float32(tileSize))
	path.LineTo(0, float32(tileSize/2))
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

func (g *Game) Draw(screen *ebiten.Image) {
	initSprites()
	screen.Fill(color.RGBA{20, 25, 30, 255})

	fs := getFontSource()

	// Título
	opShadow := &text.DrawOptions{}
	opShadow.GeoM.Translate(22, 22)
	opShadow.ColorScale.Scale(0, 0, 0, 0.8)
	text.Draw(screen, "7 COLORS", &text.GoTextFace{Source: fs, Size: 32}, opShadow)

	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 20)
	text.Draw(screen, "7 COLORS", &text.GoTextFace{Source: fs, Size: 32}, op)

	// Barra lateral
	for i, clr := range palette {
		bx := float32(30)
		by := float32(80 + i*55)
		bw := float32(45)
		bh := float32(45)

		isPlayerChoice := i == g.playerColors[0] || i == g.playerColors[1]

		vector.DrawFilledRect(screen, bx+3, by+3, bw, bh, color.RGBA{0, 0, 0, 150}, false)

		if isPlayerChoice {
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

			vector.StrokeLine(screen, bx+5, by+5, bx+bw-5, by+bh-5, 2, color.RGBA{0, 0, 0, 200}, false)
			vector.StrokeLine(screen, bx+bw-5, by+5, bx+5, by+bh-5, 2, color.RGBA{0, 0, 0, 200}, false)

			opShadowL := &text.DrawOptions{}
			opShadowL.GeoM.Translate(float64(bx+bw/4)+1, float64(by+bh/4)+1)
			opShadowL.ColorScale.Scale(0, 0, 0, 1)
			text.Draw(screen, playerText, &text.GoTextFace{Source: fs, Size: 16}, opShadowL)

			opLabel := &text.DrawOptions{}
			opLabel.GeoM.Translate(float64(bx+bw/4), float64(by+bh/4))
			r, g, b, _ := playerClr.RGBA()
			opLabel.ColorScale.Scale(float32(r)/65535, float32(g)/65535, float32(b)/65535, 1)
			text.Draw(screen, playerText, &text.GoTextFace{Source: fs, Size: 16}, opLabel)
		} else {
			vector.DrawFilledRect(screen, bx, by, bw, bh, clr, false)
			vector.StrokeLine(screen, bx, by, bx+bw, by, 2, color.RGBA{255, 255, 255, 120}, false)
			vector.StrokeLine(screen, bx, by, bx, by+bh, 2, color.RGBA{255, 255, 255, 120}, false)
			vector.StrokeLine(screen, bx+bw, by, bx+bw, by+bh, 2, color.RGBA{0, 0, 0, 120}, false)
			vector.StrokeLine(screen, bx, by+bh, bx+bw, by+bh, 2, color.RGBA{0, 0, 0, 120}, false)
		}
	}

	// Puntajes
	totalTiles := float64(boardSize * boardSize)
	p1Score := float64(len(g.playerOwned[0])) / totalTiles * 100
	p2Score := float64(len(g.playerOwned[1])) / totalTiles * 100

	drawShadowText := func(scr *ebiten.Image, txt string, size float64, fX, fY float64, clr color.Color) {
		opS := &text.DrawOptions{}
		opS.GeoM.Translate(fX+2, fY+2)
		opS.ColorScale.Scale(0, 0, 0, 0.7)
		text.Draw(scr, txt, &text.GoTextFace{Source: fs, Size: size}, opS)
		
		opT := &text.DrawOptions{}
		opT.GeoM.Translate(fX, fY)
		r, gg, b, _ := clr.RGBA()
		opT.ColorScale.Scale(float32(r)/65535, float32(gg)/65535, float32(b)/65535, 1)
		text.Draw(scr, txt, &text.GoTextFace{Source: fs, Size: size}, opT)
	}

	drawShadowText(screen, fmt.Sprintf("P1 AREA: %.1f%%", p1Score), 24, screenWidth-220, 40, color.White)
	drawShadowText(screen, fmt.Sprintf("P2 AREA: %.1f%%", p2Score), 24, screenWidth-220, 80, color.RGBA{180, 180, 180, 255})

	// Turno
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
	
	textW := float64(len(turnText)) * 11.5 // Aproximación
	drawShadowText(screen, turnText, 24, screenWidth/2-textW/2, screenHeight-50, turnColor)

	// Tablero
	offsetX := float32(screenWidth) / 2 + 50
	offsetY := float32(screenHeight) / 2 - 30

	// PASADA 1: Bases Planas
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

	// PASADA 2: Biseles 3D
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			clrIdx := g.board[y][x]
			baseClr := palette[clrIdx]

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

			if x+1 >= boardSize || g.board[y][x+1] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(dR, dG, dB, 1.0)
				screen.DrawImage(bevelBR, op)
			}
			if y+1 >= boardSize || g.board[y+1][x] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(dR, dG, dB, 1.0)
				screen.DrawImage(bevelBL, op)
			}
			if x-1 < 0 || g.board[y][x-1] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(lR, lG, lB, 1.0)
				screen.DrawImage(bevelTL, op)
			}
			if y-1 < 0 || g.board[y-1][x] != clrIdx {
				op := &ebiten.DrawImageOptions{}
				op.GeoM = drawOp.GeoM
				op.ColorScale.Scale(lR, lG, lB, 1.0)
				screen.DrawImage(bevelTR, op)
			}
		}
	}
}
