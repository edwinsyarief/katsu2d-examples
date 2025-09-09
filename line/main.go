package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/edwinsyarief/katsu2d"
)

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
	points []ebiten.Point
}

// NewGame creates a new Game object.
func NewGame() *Game {
	g := &Game{}
	g.points = make([]ebiten.Point, 0)

	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(640, 480),
		katsu2d.WithWindowTitle("Interactive Line Example"),
	)

	return g
}

func (g *Game) Update() error {
	// On left click, add a new point at the cursor position.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.points = append(g.points, ebiten.Point{X: x, Y: y})
	}

	// On right click, remove the last point.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if len(g.points) > 0 {
			g.points = g.points[:len(g.points)-1]
		}
	}

	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw all the points as small red boxes.
	for _, p := range g.points {
		vector.DrawFilledRect(screen, float32(p.X-2), float32(p.Y-2), 5, 5, color.RGBA{R: 255, A: 255}, false)
	}

	// If there are at least two points, draw lines between them.
	if len(g.points) >= 2 {
		for i := 0; i < len(g.points)-1; i++ {
			p1 := g.points[i]
			p2 := g.points[i+1]
			vector.StrokeLine(screen, float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y), 1, color.White, false)
		}
	}

	g.engine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
