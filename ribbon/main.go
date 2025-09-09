package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/edwinsyarief/katsu2d"
)

const (
	maxTrailCount = 100
)

var (
	whiteImage = ebiten.NewImage(1, 1)
)

func init() {
	whiteImage.Fill(color.White)
}

// Point represents a point in the ribbon trail.
type Point struct {
	X, Y float64
}

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
	trail  []Point
}

// NewGame creates a new Game object.
func NewGame() *Game {
	g := &Game{}
	g.trail = make([]Point, 0, maxTrailCount)

	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(640, 480),
		katsu2d.WithWindowTitle("Interactive Ribbon Trail Example"),
	)

	return g
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	g.trail = append(g.trail, Point{X: float64(x), Y: float64(y)})
	if len(g.trail) > maxTrailCount {
		g.trail = g.trail[1:]
	}
	return g.engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	if len(g.trail) < 2 {
		return
	}

	vertices := make([]ebiten.Vertex, 0, len(g.trail)*2)
	indices := make([]uint16, 0, (len(g.trail)-1)*6)

	for i := 0; i < len(g.trail); i++ {
		p1 := g.trail[i]

		var p0, p2 Point
		if i > 0 {
			p0 = g.trail[i-1]
		} else {
			p0 = p1
		}
		if i < len(g.trail)-1 {
			p2 = g.trail[i+1]
		} else {
			p2 = p1
		}

		// Get the direction of the line
		dx := p2.X - p0.X
		dy := p2.Y - p0.Y
		angle := math.Atan2(dy, dx)

		// Get the perpendicular angle
		perpAngle := angle + math.Pi/2

		// Calculate the width based on the point's age
		// The newest point is the widest, and it gets thinner.
		progress := float64(i) / float64(len(g.trail))
		width := (1 - progress) * 20 // Max width of 20 pixels

		// Create two vertices for the current point, offset by the perpendicular
		offsetX := math.Cos(perpAngle) * width / 2
		offsetY := math.Sin(perpAngle) * width / 2

		v1 := ebiten.Vertex{
			DstX: float32(p1.X + offsetX),
			DstY: float32(p1.Y + offsetY),
			SrcX: 0, SrcY: 0,
		}
		v2 := ebiten.Vertex{
			DstX: float32(p1.X - offsetX),
			DstY: float32(p1.Y - offsetY),
			SrcX: 0, SrcY: 0,
		}

		// Interpolate color from blue to transparent
		alpha := (1 - progress) * 0.8
		v1.ColorR = float32(progress) // Fades to white-ish blue
		v1.ColorG = float32(progress)
		v1.ColorB = 1.0
		v1.ColorA = float32(alpha)
		v2.ColorR = v1.ColorR
		v2.ColorG = v1.ColorG
		v2.ColorB = v1.ColorB
		v2.ColorA = v1.ColorA

		vertices = append(vertices, v1, v2)

		// Create the triangles for this segment
		if i > 0 {
			idx := uint16(i*2)
			indices = append(indices, idx-2, idx-1, idx)
			indices = append(indices, idx-1, idx+1, idx)
		}
	}

	opts := &ebiten.DrawTrianglesOptions{}
	opts.FillRule = ebiten.FillRuleFillAll
	screen.DrawTriangles(vertices, indices, whiteImage, opts)

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
