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

	entity := world.CreateEntity()
	lazyecs.AddComponent[katsu2d.TransformComponent](world, entity)
	lazyecs.AddComponent[katsu2d.ShapeComponent](world, entity)

	entity2 := world.CreateEntity()
	lazyecs.AddComponent[katsu2d.TransformComponent](world, entity2)
	lazyecs.AddComponent[katsu2d.ShapeComponent](world, entity2)

	entity3 := world.CreateEntity()
	lazyecs.AddComponent[katsu2d.TransformComponent](world, entity3)
	lazyecs.AddComponent[katsu2d.ShapeComponent](world, entity3)

	entity4 := world.CreateEntity()
	lazyecs.AddComponent[katsu2d.TransformComponent](world, entity4)
	lazyecs.AddComponent[katsu2d.ShapeComponent](world, entity4)

	t, _ := lazyecs.GetComponent[katsu2d.TransformComponent](world, entity)
	t.Init()
	t.SetPosition(ebimath.V(150, 150))
	t.SetOffset(ebimath.V2(100))

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
	shape, _ := lazyecs.GetComponent[katsu2d.ShapeComponent](world, entity)
	shape.Init(rect)

	t2, _ := lazyecs.GetComponent[katsu2d.TransformComponent](world, entity2)
	t2.Init()
	t2.SetPosition(ebimath.V(300, 300))

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
	shape2, _ := lazyecs.GetComponent[katsu2d.ShapeComponent](world, entity2)
	shape2.Init(circle)

	t3, _ := lazyecs.GetComponent[katsu2d.TransformComponent](world, entity3)
	t3.Init()
	t3.SetPosition(ebimath.V(500, 50))

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
	shape3, _ := lazyecs.GetComponent[katsu2d.ShapeComponent](world, entity3)
	shape3.Init(hexagon)

	t4, _ := lazyecs.GetComponent[katsu2d.TransformComponent](world, entity4)
	t4.Init()
	t4.SetPosition(ebimath.V(550, 400))

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
	shape4, _ := lazyecs.GetComponent[katsu2d.ShapeComponent](world, entity4)
	shape4.Init(triangle)

	g.engine.AddOverlayDrawSystem(katsu2d.NewShapeRenderSystem())

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
