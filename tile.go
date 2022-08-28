//go:build 7colors
// +build 7colors

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tile represents a space with an x,y coordinate within a Level.
type Tile struct {
	sprite *ebiten.Image
}

// AddSprite adds a sprite to the Tile.
func (t *Tile) Paint(c color.RGBA) {
	t.sprite.Fill(c)
}

// Draw draws the Tile on the screen using the provided options.
func (t *Tile) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	screen.DrawImage(t.sprite, options)
}
