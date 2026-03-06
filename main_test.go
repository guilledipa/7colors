package main

import (
	"testing"
)

func TestCapture(t *testing.T) {
	// Armamos una partida para testear
	g := NewGame()

	// Chequeamos que la captura inicial ande: los jugadores tienen que tener al menos una ficha cada uno
	if len(g.playerOwned[0]) == 0 || len(g.playerOwned[1]) == 0 {
		t.Errorf("La captura inicial falló: El Jugador 1 tiene %d fichas, el Jugador 2 tiene %d",
			len(g.playerOwned[0]), len(g.playerOwned[1]))
	}

	// Nos aseguramos que no arranquen con el mismo color
	if g.playerColors[0] == g.playerColors[1] {
		t.Errorf("Los colores de entrada deberían ser distintos, pero los dos tienen el %d", g.playerColors[0])
	}
}
