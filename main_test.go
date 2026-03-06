package main

import (
	"testing"
)

func TestCapture(t *testing.T) {
	g := NewGame()

	// Test basic capture logic - ensure players own at least one tile
	if len(g.playerOwned[0]) == 0 || len(g.playerOwned[1]) == 0 {
		t.Errorf("Initial capture failed: Player 1 owned %d tiles, Player 2 owned %d tiles",
			len(g.playerOwned[0]), len(g.playerOwned[1]))
	}

	// Ensure player colors are different
	if g.playerColors[0] == g.playerColors[1] {
		t.Errorf("Initial colors should be different, but both are %d", g.playerColors[0])
	}
}
