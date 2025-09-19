package main

import (
	"fmt"
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	WindowWidth  = 640
	WindowHeight = 360
)

var ParticleTypes = []string{"Fire", "Rain", "Whimsical"}
var SelectedParticle = 0

type DebugSystem struct {
	debugImg *ebiten.Image
}

func (self *DebugSystem) Draw(world *katsu2d.World, renderer *katsu2d.BatchRenderer) {
	screen := renderer.GetScreen()

	if self.debugImg == nil {
		self.debugImg = ebiten.NewImage(WindowWidth/2, WindowHeight/2)
	}
	self.debugImg.Fill(color.Transparent)
	ebitenutil.DebugPrintAt(self.debugImg,
		fmt.Sprintf("FPS: %2.f\nPress [Q]/[E] to change particle type\nSelected Particle: %s", ebiten.ActualFPS(), ParticleTypes[SelectedParticle]),
		10, 10)

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(2, 2)
	screen.DrawImage(self.debugImg, &opts)
}

type MainSystem struct {
	selectedParticle int
	prevParticle     int
}

func (self *MainSystem) Update(world *katsu2d.World, dt float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		self.selectedParticle++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		self.selectedParticle--
	}

	self.selectedParticle = ebimath.Clamp(self.selectedParticle, 0, 2)
	SelectedParticle = self.selectedParticle

	if self.selectedParticle != self.prevParticle {
		self.prevParticle = self.selectedParticle

		emitterEntity := world.Query(katsu2d.CTParticleEmitter)[0]
		world.RemoveComponent(emitterEntity, katsu2d.CTParticleEmitter)

		t, _ := world.GetComponent(emitterEntity, katsu2d.CTTransform)
		transform := t.(*katsu2d.TransformComponent)

		switch self.selectedParticle {
		case 1:
			rainPreset := katsu2d.RainPreset(0)
			rainPreset.InitialColorMin = color.RGBA{255, 255, 255, 255}
			rainPreset.InitialColorMax = color.RGBA{255, 255, 255, 255}
			rainPreset.TargetColorMin = color.RGBA{255, 255, 255, 0}
			rainPreset.TargetColorMax = color.RGBA{255, 255, 255, 0}
			rainPreset.MinScale = 2
			rainPreset.MaxScale = 3
			rainPreset.InitialParticleSpeedMin = 0
			rainPreset.InitialParticleSpeedMax = 0
			rainPreset.DirectionMode = katsu2d.ParticleDirectionModeLinear
			rainPreset.Gravity = ebimath.V(0, 900)
			rainPreset.ParticleSpawnOffset = ebimath.V(320, 0)
			world.AddComponent(emitterEntity, rainPreset)
			transform.SetPosition(ebimath.V(0, -10))
		case 2:
			whimsicalPreset := katsu2d.WhimsicalPreset(0)
			whimsicalPreset.MinScale = 5
			whimsicalPreset.MaxScale = 10
			world.AddComponent(emitterEntity, whimsicalPreset)
			transform.SetPosition(ebimath.V(320/2, 180/2))
		default:
			// --- Particle Emitter Setup ---
			firePreset := katsu2d.FirePreset(0)
			firePreset.MinScale = 10
			firePreset.MaxScale = 10.5
			firePreset.ParticleSpawnOffset = ebimath.V2(20)
			firePreset.ScaleMode = katsu2d.ParticleScaleModeScaleInOut
			firePreset.Gravity = ebimath.V(0, -150)
			world.AddComponent(emitterEntity, firePreset)
			transform.SetPosition(ebimath.V(320/2, 180/2))
		}
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
		katsu2d.WithWindowSize(WindowWidth, WindowHeight),
		katsu2d.WithWindowTitle("Particle Example"),
	)

	tm := g.engine.TextureManager()
	world := g.engine.World()

	// --- Player Entity Setup ---
	playerEntity := world.CreateEntity()
	playerTransform := katsu2d.NewTransformComponent()
	// Start the emitter at the bottom-center of the screen
	playerTransform.SetPosition(ebimath.V(WindowWidth/4, WindowHeight/4))
	world.AddComponent(playerEntity, playerTransform)

	// --- Particle Emitter Setup ---
	firePreset := katsu2d.FirePreset(0)
	firePreset.MinScale = 10
	firePreset.MaxScale = 10.5
	firePreset.ParticleSpawnOffset = ebimath.V2(20)
	firePreset.ScaleMode = katsu2d.ParticleScaleModeScaleInOut
	firePreset.Gravity = ebimath.V(0, -150)
	world.AddComponent(playerEntity, firePreset)

	renderer := katsu2d.NewLayerSystem(
		world,
		WindowWidth/2, WindowHeight/2,
		katsu2d.AddSystem(katsu2d.NewParticleRenderSystem(tm)))

	// --- System Setup ---
	g.engine.AddUpdateSystem(&MainSystem{})
	g.engine.AddUpdateSystem(katsu2d.NewParticleEmitterSystem(tm))
	g.engine.AddUpdateSystem(katsu2d.NewParticleUpdateSystem())
	g.engine.AddBackgroundDrawSystem(renderer)
	g.engine.AddOverlayDrawSystem(&DebugSystem{})

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
