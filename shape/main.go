package main

import (
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
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
		katsu2d.WithWindowSize(640, 480),
		katsu2d.WithWindowTitle("Shape Example"),
		katsu2d.WithUpdateSystem(katsu2d.NewInputSystem()),
	)

	world := g.engine.World()

	entity := world.CreateEntity()
	t := katsu2d.NewTransformComponent()
	t.SetPosition(ebimath.V(150, 150))
	t.SetOffset(ebimath.V2(100))
	world.AddComponent(entity, t)

	r := katsu2d.NewRectangleComponent(200, 200, color.RGBA{255, 255, 0, 255})
	r.SetCornerRadius(20, 0, 0, 20)
	r.SetStroke(20, color.RGBA{255, 0, 0, 255})
	r.SetColor(
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	r.SetStrokeColor(
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	world.AddComponent(entity, r)

	entity2 := world.CreateEntity()
	t2 := katsu2d.NewTransformComponent()
	t2.SetPosition(ebimath.V(300, 300))
	world.AddComponent(entity2, t2)

	circle := katsu2d.NewCircleComponent(50, color.RGBA{255, 255, 0, 255})
	circle.SetStroke(20, color.RGBA{255, 0, 0, 255})
	circle.SetColor(
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	circle.SetStrokeColor(
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	world.AddComponent(entity2, circle)

	entity3 := world.CreateEntity()
	t3 := katsu2d.NewTransformComponent()
	t3.SetPosition(ebimath.V(500, 50))
	world.AddComponent(entity3, t3)

	hexagon := katsu2d.NewHexagonComponent(100, color.RGBA{255, 255, 0, 255})
	hexagon.SetCornerRadius(10)
	hexagon.SetStroke(20, color.RGBA{255, 0, 0, 255})
	hexagon.SetColor(
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	hexagon.SetStrokeColor(
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	world.AddComponent(entity3, hexagon)

	entity4 := world.CreateEntity()
	t4 := katsu2d.NewTransformComponent()
	t4.SetPosition(ebimath.V(550, 400))
	world.AddComponent(entity4, t4)

	triangle := katsu2d.NewTriangleComponent(100, 100, color.RGBA{255, 255, 0, 255})
	triangle.SetCornerRadius(10)
	triangle.SetStroke(20, color.RGBA{255, 0, 0, 255})
	triangle.SetColor(
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	triangle.SetStrokeColor(
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
	)
	world.AddComponent(entity4, triangle)

	g.engine.AddOverlayDrawSystem(katsu2d.NewShapeRenderSystem())

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
