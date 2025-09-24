package main

import (
	_ "image/png"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
}

var rand *ebimath.Rand

// NewGame creates a new Game object and sets up the engine.
func NewGame() *Game {
	rand = ebimath.Random()

	g := &Game{}

	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(640, 480),
		katsu2d.WithWindowTitle("Foliage Example"),
		katsu2d.WithVsyncEnabled(true),
	)

	tm := g.engine.TextureManager()
	world := g.engine.World()

	// --- Texture Loading ---
	foliageImage, _, _ := ebitenutil.NewImageFromFile("./grass.png")
	foliageTextureID := tm.Add(foliageImage)

	// --- System Setup ---
	g.engine.AddUpdateSystem(katsu2d.NewFoliageSystem())
	g.engine.AddBackgroundDrawSystem(katsu2d.NewSpriteRenderSystem(tm))

	// --- Foliage Controller ---
	foliageControllerEntity := world.CreateEntity()
	foliageController, _ := lazyecs.AddComponent[katsu2d.FoliageControllerComponent](world, foliageControllerEntity)
	foliageController.Init(
		katsu2d.WithFoliageWindForce(200), // Radians
		katsu2d.WithFoliageWindSpeed(3.95),
		katsu2d.WithFoliageRippleStrength(50),
	)

	createFoliage(world, tm, foliageTextureID, 40, 120, ebimath.V(0.5, 1.0))

	return g
}

func createFoliage(world *lazyecs.World, tm *katsu2d.TextureManager, textureID int, x, y float64, pivot ebimath.Vector) {
	entity := world.CreateEntity()

	// Transform
	transform, _ := lazyecs.AddComponent[katsu2d.TransformComponent](world, entity)
	transform.Init()
	transform.SetPosition(ebimath.V(x, y))

	// Sprite with a grid mesh
	img := tm.Get(textureID)
	sprite, _ := lazyecs.AddComponent[katsu2d.SpriteComponent](world, entity)
	sprite.Init(textureID, img.Bounds())
	sprite.DstW = 512
	sprite.DstH = 512
	sprite.SetGrid(5, 5)

	// Foliage
	foliage, _ := lazyecs.AddComponent[katsu2d.FoliageComponent](world, entity)
	foliage.TextureID = textureID
	foliage.SwaySeed = rand.FloatRange(0, 100)
	foliage.Pivot = pivot
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
