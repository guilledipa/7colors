package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
)

var (
	colors = []color.RGBA{
		{0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0x0, 0xff},
		{0xff, 0x0, 0xff, 0xff},
		{0xff, 0x0, 0x0, 0xff},
		{0x0, 0xff, 0xff, 0xff},
		{0x0, 0xff, 0x0, 0xff},
		{0x0, 0x0, 0xff, 0xff},
		{0x0, 0x0, 0x0, 0xff},
	}
)

// Level represents a Game level.
type Level struct {
	w, h int

	tiles    [][]*Tile // (Y,X) array of tiles
	tileSize int
}

// Tile returns the tile at the provided coordinates, or nil.
func (l *Level) Tile(x, y int) *Tile {
	if x >= 0 && y >= 0 && x < l.w && y < l.h {
		return l.tiles[y][x]
	}
	return nil
}

// Size returns the size of the Level.
func (l *Level) Size() (width, height int) {
	return l.w, l.h
}

// NewLevel returns a new randomly generated Level.
func NewLevel() (*Level, error) {
	// Create a 108x108 Level.
	l := &Level{
		w:        108,
		h:        108,
		tileSize: 64,
	}

	// Load embedded SpriteSheet.
	ss, err := LoadSpriteSheet(l.tileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded spritesheet: %s", err)
	}

	// Generate a unique permutation each time.
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// Fill each tile with one or more sprites randomly.
	l.tiles = make([][]*Tile, l.h)
	for y := 0; y < l.h; y++ {
		l.tiles[y] = make([]*Tile, l.w)
		for x := 0; x < l.w; x++ {
			t := &Tile{}
			isBorderSpace := x == 0 || y == 0 || x == l.w-1 || y == l.h-1
			val := r.Intn(1000)
			switch {
			case isBorderSpace || val < 275:
				t.Paint(colors[1])
			case val < 285:
				t.Paint(colors[2])
			case val < 288:
				t.Paint(colors[3])
			case val < 289:
				t.Paint(colors[4])
				t.Paint(colors[5])
			case val < 290:
				t.Paint(colors[6])
			default:
				t.Paint(colors[7])
			}
			l.tiles[y][x] = t
		}
	}
	return l, nil
}
