package main

import (
	"fmt"
	"log"

	_ "image/png"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GrassSystem struct {
	debugImg *ebiten.Image
}

func (self *GrassSystem) Update(world *lazyecs.World, dt float64) {
	query := world.Query(katsu2d.CTGrassController)
	var grassCtrl *katsu2d.GrassControllerComponent
	for query.Next() {
		for _, entity := range query.Entities() {
			grassCtrl, _ = lazyecs.GetComponent[katsu2d.GrassControllerComponent](world, entity)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		grassCtrl.AddStrongWindGust(katsu2d.StrongWindGust{
			Width:           200,
			StartPos:        ebimath.V(-100, 0),
			EndPos:          ebimath.V2(500),
			Strength:        500,
			Length:          200,
			Duration:        3.5,
			FadeInDuration:  .25,
			FadeOutDuration: .75,
		})
	}

	x, y := ebiten.CursorPosition()
	grassCtrl.SetForcePositions(katsu2d.ForceSource{
		Radius:   100,
		Position: ebimath.V(float64(x), float64(y)),
		Strength: 100,
	})
}

func (self *GrassSystem) Draw(world *lazyecs.World, renderer *katsu2d.BatchRenderer) {
	screen := renderer.GetScreen()
	renderer.Flush()

	bladeAmount := 0
	query := world.Query(katsu2d.CTGrass)
	for query.Next() {
		bladeAmount += query.Count()
	}

	if self.debugImg == nil {
		self.debugImg = ebiten.NewImage(320, 180)
	}
	self.debugImg.Clear()
	ebitenutil.DebugPrintAt(self.debugImg, fmt.Sprintf("FPS: %.2f\nBlade Amount: %d", ebiten.ActualFPS(), bladeAmount), 5, 5)

	ops := ebiten.DrawImageOptions{}
	ops.GeoM.Scale(2, 2)
	screen.DrawImage(self.debugImg, &ops)
}

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
		katsu2d.WithWindowTitle("Grass Example"),
		katsu2d.WithUpdateSystem(katsu2d.NewGrassControllerSystem()),
		katsu2d.WithWindowResizeMode(ebiten.WindowResizingModeEnabled),
	)

	tm := g.engine.TextureManager()
	world := g.engine.World()

	// --- Texture Loading ---
	// Grass
	grass, _, _ := ebitenutil.NewImageFromFile("./wheat.png")
	texId := tm.Add(grass) // ID 1: "grass"

	// grass controller
	entity := world.CreateEntity()
	transform, _ := lazyecs.AddComponent[katsu2d.TransformComponent](world, entity)
	transform.Init()
	grassController, _ := lazyecs.AddComponent[katsu2d.GrassControllerComponent](world, entity)
	grassController.Init(world, tm,
		640, 480, texId, transform.Z,
		katsu2d.WithGrassOrderable(true),
		katsu2d.WithGrassDensity(10),
		katsu2d.WithGrassWindDirection(1, 0),
		katsu2d.WithGrassWindForce(0.5),
		katsu2d.WithGrassWindSpeed(3.5),
		katsu2d.WithGrassAreas([]katsu2d.Area{
			{X1: 0, Y1: 0, X2: 20, Y2: 20},
		}),
	)

	ls := katsu2d.NewLayerSystem( 640, 480,
		katsu2d.AddSystem(katsu2d.NewOrderableSystem(tm)),
	)

	g.engine.AddBackgroundDrawSystem(ls)
	g.engine.AddOverlayDrawSystem(&GrassSystem{})

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
