package main

import (
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/katsu2d/line"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type RibbonSystem struct {
	ribbon  *line.RibbonTrails
	isDebug bool
}

func (self *RibbonSystem) Update(world *katsu2d.World, dt float64) {
	x, y := ebiten.CursorPosition()
	self.ribbon.AddPoint(ebimath.V(float64(x), float64(y)))
	self.ribbon.Update(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		self.isDebug = !self.isDebug
		self.ribbon.SetDebugDraw(self.isDebug)
	}
}

func (self *RibbonSystem) Draw(world *katsu2d.World, renderer *katsu2d.BatchRenderer) {
	screen := renderer.GetScreen()
	topts := &ebiten.DrawTrianglesOptions{}
	topts.Blend.BlendOperationRGB = ebiten.BlendOperationMax
	topts.AntiAlias = true
	self.ribbon.Draw(screen, topts)
	ebitenutil.DebugPrintAt(screen, "Move your mouse to draw trails\nPress [Space] to toggle debug drawing", 10, 10)
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
		katsu2d.WithWindowTitle("Interactive Ribbon Trail Example"),
	)

	ribbon := line.NewRibbonTrails()
	ribbon.SetColors(color.RGBA{R: 255, G: 255, B: 255, A: 255}, color.RGBA{R: 255, G: 255, B: 255, A: 0})
	ribbon.SetJointMode(line.LineJointRound)
	ribbon.SetLifetime(0.75)
	ribbon.SetWidths(20, 0)

	g.engine.AddOverlayDrawSystem(&RibbonSystem{
		ribbon: ribbon,
	})

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
