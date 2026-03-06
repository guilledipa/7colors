# 7 Colors - Ebitengine Clone

A faithful recreation of the 1991 DOS puzzle game **7 Colors** (originally by Gamos/Infogrames), implemented in Go using the **Ebitengine** framework.

## Game Rules

7 Colors is a turn-based board game played on a grid of diamonds of 7 different colors.

- **Objective:** Control more than 50% of the board tiles.
- **Starting Points:** Player 1 starts at the bottom-left corner, and Player 2 starts at the top-right corner.
- **How to Play:** On your turn, choose a color from the sidebar. All tiles you currently own will change to that color, and you will capture all adjacent tiles that share the new color.
- **Restrictions:** You cannot choose the color currently owned by yourself or your opponent.

## Features

- **DOS Aesthetic:** Classic EGA-like color palette and diamond-shaped grid.
- **Score Tracking:** Real-time percentage of the board controlled by each player.
- **Responsive Controls:** Select colors using mouse clicks on the sidebar.
- **Game End:** Automatic detection of the winner and a simple restart mechanism.

## How to Run

Ensure you have [Go](https://go.dev/) installed on your system.

1. Clone the repository.
2. Run the game:
   ```bash
   go run .
   ```

## Controls

- **Mouse Left Click:** Select a color from the left sidebar.
- **Restart:** Click anywhere after the game ends to start a new match.

## Framework

This game is built with [Ebitengine](https://ebitengine.org/), an open-source game engine for Go.
