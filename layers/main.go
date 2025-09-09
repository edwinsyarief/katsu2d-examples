package main

import (
	"log"

	"github.com/edwinsyarief/katsu2d"
)

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
}

// NewGame creates a new Game object and sets up the engine.
func NewGame() *Game {
	g := &Game{}

	// --- Engine Setup ---
	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(320, 240),
		katsu2d.WithWindowTitle("Layers Example"),
		katsu2d.WithUpdateSystem(katsu2d.NewInputSystem()),
	)

	// TODO

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
