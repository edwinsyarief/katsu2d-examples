package main

import (
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
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

func (self *PlayerSystem) Update(world *katsu2d.World, dt float64) {
	// Find the player entity using its tag.
	var player katsu2d.Entity
	found := false
	for _, e := range world.Query(katsu2d.CTTag, katsu2d.CTTransform, katsu2d.CTInput) {
		tag, _ := world.GetComponent(e, katsu2d.CTTag)
		if tag.(*katsu2d.TagComponent).Tag == PlayerTag {
			player = e
			found = true
			break
		}
	}

	if !found {
		return
	}

	t, _ := world.GetComponent(player, katsu2d.CTTransform)
	transform := t.(*katsu2d.TransformComponent)

	i, _ := world.GetComponent(player, katsu2d.CTInput)
	input := i.(*katsu2d.InputComponent)

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

	// --- Engine Setup ---
	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(320, 240),
		katsu2d.WithWindowTitle("Basic Example"),
		katsu2d.WithUpdateSystem(katsu2d.NewInputSystem()),
	)

	tm := g.engine.TextureManager()
	world := g.engine.World()

	playerImg := ebiten.NewImage(16, 16)
	playerImg.Fill(color.White)
	playerTexID := tm.Add(playerImg) // ID 4: Player

	// --- Player Setup ---
	playerEntity := world.CreateEntity()
	playerTransform := katsu2d.NewTransformComponent()
	playerTransform.SetPosition(ebimath.V(80, 60))
	playerTransform.Z = 1 // Set player's Z to the same layer as the upper grid tiles
	world.AddComponent(playerEntity, playerTransform)
	world.AddComponent(playerEntity, katsu2d.NewSpriteComponent(playerTexID, playerImg.Bounds()))
	world.AddComponent(playerEntity, katsu2d.NewTagComponent(PlayerTag))
	world.AddComponent(playerEntity, katsu2d.NewInputComponent(keybindings))

	// --- System Setup ---
	// The order of these systems is important for this rendering technique.
	g.engine.AddUpdateSystem(&PlayerSystem{})

	g.engine.AddBackgroundDrawSystem(katsu2d.NewSpriteRenderSystem(world, tm))

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
