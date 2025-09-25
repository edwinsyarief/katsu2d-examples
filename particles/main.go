package main

import (
	"fmt"
	"image/color"
	"log"

	ebimath "github.com/edwinsyarief/ebi-math"
	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
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

func (self *DebugSystem) Draw(world *lazyecs.World, renderer *katsu2d.BatchRenderer) {
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

func (self *MainSystem) Update(world *lazyecs.World, dt float64) {
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

		query := world.Query(katsu2d.CTParticleEmitter)
		for query.Next() {
			for _, entity := range query.Entities() {
				lazyecs.RemoveComponent[katsu2d.ParticleEmitterComponent](world, entity)

				transform, _ := lazyecs.GetComponent[katsu2d.TransformComponent](world, entity)
				var newPreset *katsu2d.ParticleEmitterComponent

				switch self.selectedParticle {
				case 1:
					newPreset = katsu2d.RainPreset(rainTextId)
					newPreset.MaxParticles = 5000
					newPreset.EmitRate = 2000
					newPreset.InitialColorMin = color.RGBA{255, 255, 255, 255}
					newPreset.InitialColorMax = color.RGBA{255, 255, 255, 255}
					newPreset.TargetColorMin = color.RGBA{255, 255, 255, 0}
					newPreset.TargetColorMax = color.RGBA{255, 255, 255, 0}
					newPreset.MinScale = 2
					newPreset.MaxScale = 3
					newPreset.InitialParticleSpeedMin = 0
					newPreset.InitialParticleSpeedMax = 50
					newPreset.DirectionMode = katsu2d.ParticleDirectionModeLinear
					newPreset.Gravity = ebimath.V(0, 900)
					newPreset.ParticleSpawnOffset = ebimath.V(320, 0)
					transform.SetPosition(ebimath.V(0, -50))
				case 2:
					newPreset = katsu2d.WhimsicalPreset(0)
					newPreset.MinScale = 5
					newPreset.MaxScale = 10
					transform.SetPosition(ebimath.V(320/2, 180/2))
				default:
					// --- Particle Emitter Setup ---
					newPreset = katsu2d.FirePreset(0)
					newPreset.MinScale = 10
					newPreset.MaxScale = 10.5
					newPreset.ParticleSpawnOffset = ebimath.V2(20)
					newPreset.ScaleMode = katsu2d.ParticleScaleModeScaleInOut
					newPreset.Gravity = ebimath.V(0, -150)
					transform.SetPosition(ebimath.V(320/2, 180/2))
				}

				lazyecs.SetComponent(world, entity, *newPreset)
			}
		}
	}
}

// Game implements ebiten.Game interface.
type Game struct {
	engine *katsu2d.Engine
}

var rainTextId = 0

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

	rainImage := ebiten.NewImage(1, 5)
	rainImage.Fill(color.White)
	rainTextId = tm.Add(rainImage)

	// --- Player Entity Setup ---
	playerEntity := world.CreateEntity()
	playerTransform := katsu2d.NewTransformComponent()
	// Start the emitter at the bottom-center of the screen
	playerTransform.SetPosition(ebimath.V(WindowWidth/4, WindowHeight/4))
	lazyecs.SetComponent(world, playerEntity, *playerTransform)

	// --- Particle Emitter Setup ---
	firePreset := katsu2d.FirePreset(0)
	firePreset.MinScale = 10
	firePreset.MaxScale = 10.5
	firePreset.ParticleSpawnOffset = ebimath.V2(20)
	firePreset.ScaleMode = katsu2d.ParticleScaleModeScaleInOut
	firePreset.Gravity = ebimath.V(0, -150)
	lazyecs.SetComponent(world, playerEntity, *firePreset)

	renderer := katsu2d.NewLayerSystem(
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
