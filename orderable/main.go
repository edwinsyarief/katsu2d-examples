package main

import (
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ActionMoveUp    katsu2d.Action = "move_up"
	ActionMoveDown  katsu2d.Action = "move_down"
	ActionMoveLeft  katsu2d.Action = "move_left"
	ActionMoveRight katsu2d.Action = "move_right"
)

var keybindings = map[katsu2d.Action][]katsu2d.KeyConfig{
	ActionMoveUp:    {{Primary: ebiten.KeyW}, {Primary: ebiten.KeyUp}},
	ActionMoveDown:  {{Primary: ebiten.KeyS}, {Primary: ebiten.KeyDown}},
	ActionMoveLeft:  {{Primary: ebiten.KeyA}, {Primary: ebiten.KeyLeft}},
	ActionMoveRight: {{Primary: ebiten.KeyD}, {Primary: ebiten.KeyRight}},
}

const PlayerTag = "player"

// PlayerSystem is a simple system to move the player.
type PlayerSystem struct{}

func (self *PlayerSystem) Update(world *lazyecs.World, dt float64) {
	// Find the player entity using its tag.
	found := false
	query := world.Query(katsu2d.CTTag)
	var transform *katsu2d.TransformComponent
	var input katsu2d.InputComponent
	for query.Next() {
		tags, _ := lazyecs.GetComponentSlice[katsu2d.TagComponent](query)
		inputs, _ := lazyecs.GetComponentSlice[katsu2d.InputComponent](query)
		for i, entity := range query.Entities() {
			if tags[i].Tag == PlayerTag {
				found = true
				transform, _ = lazyecs.GetComponent[katsu2d.TransformComponent](world, entity)
				input = inputs[i]
				break
			}
		}
	}

	if !found {
		return
	}

	speed := 60.0 // pixels per second
	var velocity ebimath.Vector
	if input.IsPressed(ActionMoveUp) {
		velocity.Y = -1
	}
	if input.IsPressed(ActionMoveDown) {
		velocity.Y = 1
	}
	if input.IsPressed(ActionMoveLeft) {
		velocity.X = -1
	}
	if input.IsPressed(ActionMoveRight) {
		velocity.X = 1
	}

	if !velocity.IsZero() {
		transform.SetPosition(transform.Position().Add(velocity.Normalize().ScaleF(speed * dt)))
	}
}

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
}

// NewGame creates a new Game object and sets up the engine.
func NewGame() *Game {
	g := &Game{}

	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(600, 480),
		katsu2d.WithWindowTitle("Render Order Example"),
		katsu2d.WithClearScreenEachFrame(false),
	)

	tm := g.engine.TextureManager()
	world := g.engine.World()

	// --- Texture Loading ---
	treeImg := ebiten.NewImage(25, 50)
	treeImg.Fill(color.RGBA{R: 93, G: 62, B: 4, A: 255})
	treeTexID := tm.Add(treeImg)

	playerImg := ebiten.NewImage(25, 25)
	playerImg.Fill(color.White)
	playerTexID := tm.Add(playerImg)

	particleImg := ebiten.NewImage(6, 6)
	particleImg.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	particleTexID := tm.Add(particleImg)

	// --- Entities ---

	// Create some trees
	for i := 0; i < 10; i++ {
		treeEntity := world.CreateEntity()
		treeTransform := katsu2d.NewTransformComponent()
		treeTransform.SetPosition(ebimath.V(float64(i*40+20), 100))
		lazyecs.SetComponent(world, treeEntity, *treeTransform)

		treeSprite := katsu2d.NewSpriteComponent(treeTexID, treeImg.Bounds())
		lazyecs.SetComponent(world, treeEntity, *treeSprite)

		orderable := katsu2d.NewOrderableComponent(func() float64 {
			return treeTransform.Position().Y + float64(treeSprite.DstH)
		})
		lazyecs.SetComponent(world, treeEntity, *orderable)
	}

	for i := 0; i < 10; i++ {
		treeEntity := world.CreateEntity()
		treeTransform := katsu2d.NewTransformComponent()
		treeTransform.SetPosition(ebimath.V(float64(i*40+20), 200))
		lazyecs.SetComponent(world, treeEntity, *treeTransform)

		treeSprite := katsu2d.NewSpriteComponent(treeTexID, treeImg.Bounds())
		lazyecs.SetComponent(world, treeEntity, *treeSprite)

		orderable := katsu2d.NewOrderableComponent(func() float64 {
			return treeTransform.Position().Y + float64(treeSprite.DstH)
		})
		lazyecs.SetComponent(world, treeEntity, *orderable)
	}

	// Player
	playerEntity := world.CreateEntity()

	playerTransform := katsu2d.NewTransformComponent()
	playerTransform.SetPosition(ebimath.V(160, 120))
	lazyecs.SetComponent(world, playerEntity, *playerTransform)

	playerSprite := katsu2d.NewSpriteComponent(playerTexID, playerImg.Bounds())
	lazyecs.SetComponent(world, playerEntity, *playerSprite)

	playerTag := katsu2d.NewTagComponent(PlayerTag)
	lazyecs.SetComponent(world, playerEntity, *playerTag)

	playerInput := katsu2d.NewInputComponent(keybindings)
	lazyecs.SetComponent(world, playerEntity, *playerInput)

	orderable := katsu2d.NewOrderableComponent(func() float64 {
		return playerTransform.Position().Y + float64(playerSprite.DstH)
	})
	lazyecs.SetComponent(world, playerEntity, *orderable)

	// Particle Emitter
	fireEmitter := katsu2d.FirePreset(particleTexID)
	lazyecs.SetComponent(world, playerEntity, *fireEmitter)

	// --- Systems ---
	g.engine.AddUpdateSystem(katsu2d.NewInputSystem())
	g.engine.AddUpdateSystem(&PlayerSystem{})
	g.engine.AddUpdateSystem(katsu2d.NewParticleEmitterSystem(tm))
	g.engine.AddUpdateSystem(katsu2d.NewParticleUpdateSystem())
	g.engine.AddOverlayDrawSystem(katsu2d.NewOrderableSystem(tm))

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
