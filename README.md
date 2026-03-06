# 7 Colors - Clon de Ebitengine

Una recreación fiel del juego de puzles de DOS de 1991 **7 Colors** (originalmente de Gamos/Infogrames), implementado en Go usando el framework **Ebitengine**.

## Reglas del Juego

7 Colors es un juego de tablero por turnos que se juega en una cuadrícula de diamantes de 7 colores diferentes.

- **Objetivo:** Controlar más del 50% de las fichas del tablero.
- **Puntos de Inicio:** El Jugador 1 arranca en la esquina inferior izquierda y el Jugador 2 arranca en la esquina superior derecha.
- **Cómo Jugar:** En tu turno, elegí un color de la barra lateral. Todas las fichas que ya tenés van a cambiar a ese color y vas a capturar todas las fichas de al lado que tengan el mismo color.
- **Restricciones:** No podés elegir el color que ya tenés ni el que tiene tu oponente.

## Características

- **Estética DOS:** Paleta de colores clásica tipo EGA y cuadrícula de diamantes.
- **Seguimiento de Puntaje:** Porcentaje del tablero controlado por cada jugador en tiempo real.
- **Controles Piolas:** Seleccioná los colores haciendo clic con el mouse en la barra lateral.
- **Fin del Juego:** Detección automática del ganador y un mecanismo simple para reiniciar.

## Cómo Ejecutarlo

Asegurate de tener [Go](https://go.dev/) instalado en tu sistema.

1. Cloná el repositorio.
2. Corré el juego:
   ```bash
   go run .
   ```

## Controles

- **Clic Izquierdo del Mouse:** Elegí un color de la barra lateral izquierda.
- **Reiniciar:** Hacé clic en cualquier lado después de que termine la partida para empezar de nuevo.

## Framework

Este juego está hecho con [Ebitengine](https://ebitengine.org/), un motor de juegos de código abierto para Go.
