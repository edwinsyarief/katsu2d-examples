package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/katsu2d/line"
)

type LineSystem struct {
	line          *line.Line
	debug, closed bool
}

func (self *LineSystem) Update(world *katsu2d.World, dt float64) {
	// On left click, add a new point at the cursor position.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		self.line.AddPoint(ebimath.V(float64(x), float64(y)))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		self.closed = !self.closed
		self.line.SetIsClosed(self.closed)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		self.line.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		self.debug = !self.debug
		self.line.SetDebugDraw(self.debug)
	}
}

func (self *LineSystem) Draw(world *katsu2d.World, renderer *katsu2d.BatchRenderer) {
	screen := renderer.GetScreen()

	op := ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	self.line.Draw(screen, &op)

	ebitenutil.DebugPrintAt(screen,
		"Click anywhere to add point to draw line\nPress [Space] to toggle closed line\nPress [R] to reset\nPress [D] to toggle debug draw", 10, 10)
}

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
}

// NewGame creates a new Game object.
func NewGame() *Game {
	g := &Game{}

	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(640, 480),
		katsu2d.WithWindowTitle("Interactive Line Example"),
	)

	l := line.NewLine()
	l.SetWidth(20)
	l.SetJointMode(line.LineJointRound)
	g.engine.AddOverlayDrawSystem(&LineSystem{
		line: l,
	})

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
