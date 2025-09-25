package main

import (
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
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

	entities := world.CreateEntities(4)

	// Shape 1
	t := katsu2d.NewTransformComponent()
	t.SetPosition(ebimath.V(150, 150))
	t.SetOffset(ebimath.V2(100))
	lazyecs.SetComponent(world, entities[0], *t)

	rect := katsu2d.NewRectangleShape(200, 200, color.RGBA{255, 255, 0, 255})
	rect.SetCornerRadius(20, 0, 0, 20)
	rect.SetStroke(20, color.RGBA{255, 0, 0, 255})
	rect.SetColor(
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	rect.SetStrokeColor(
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 0, 0, 255},
	)
	shape := katsu2d.NewShapeComponent(rect)
	lazyecs.SetComponent(world, entities[0], *shape)

	t2 := katsu2d.NewTransformComponent()
	t2.SetPosition(ebimath.V(300, 300))
	lazyecs.SetComponent(world, entities[1], *t2)

	circle := katsu2d.NewCircleShape(50, color.RGBA{255, 255, 0, 255})
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
	shape2 := katsu2d.NewShapeComponent(circle)
	lazyecs.SetComponent(world, entities[1], *shape2)

	t3 := katsu2d.NewTransformComponent()
	t3.SetPosition(ebimath.V(500, 50))
	lazyecs.SetComponent(world, entities[2], *t3)

	hexagon := katsu2d.NewHexagonShape(100, color.RGBA{255, 255, 0, 255})
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
	shape3 := katsu2d.NewShapeComponent(hexagon)
	lazyecs.SetComponent(world, entities[2], *shape3)

	t4 := katsu2d.NewTransformComponent()
	t4.SetPosition(ebimath.V(550, 400))
	lazyecs.SetComponent(world, entities[3], *t4)

	triangle := katsu2d.NewTriangleShape(100, 100, color.RGBA{255, 255, 0, 255})
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
	shape4 := katsu2d.NewShapeComponent(triangle)
	lazyecs.SetComponent(world, entities[3], *shape4)

	g.engine.AddOverlayDrawSystem(katsu2d.NewShapeRenderSystem())

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
