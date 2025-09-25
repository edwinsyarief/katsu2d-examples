package main

import (
	"log"

	"github.com/edwinsyarief/katsu2d"
	"github.com/edwinsyarief/lazyecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ActionPlayMusic katsu2d.Action = "play_music"
	ActionPlaySfx   katsu2d.Action = "play_sfx"
)

var keybindings = map[katsu2d.Action][]katsu2d.KeyConfig{
	ActionPlayMusic: {{Primary: ebiten.KeySpace}},
	ActionPlaySfx:   {{Primary: ebiten.KeyS}},
}

// AudioSystem is a simple system to play a sound when the spacebar is pressed.
type AudioSystem struct {
	audioManager  *katsu2d.AudioManager
	trackID       katsu2d.TrackID
	playbackID    katsu2d.PlaybackID
	stackConfig   *katsu2d.StackingConfig
	sfxID         katsu2d.TrackID
	sfxPlaybackID katsu2d.PlaybackID
}

func (self *AudioSystem) Update(world *lazyecs.World, dt float64) {
	query := world.Query(katsu2d.CTInput)
	for query.Next() {
		inputs, _ := lazyecs.GetComponentSlice[katsu2d.InputComponent](query)
		for _, input := range inputs {
			if input.IsJustPressed(ActionPlayMusic) {
				self.playbackID, _ = self.audioManager.FadeSound(self.trackID, 0, 3.75, katsu2d.AudioFadeIn, self.stackConfig)
			}
			if input.IsJustPressed(ActionPlaySfx) {
				self.sfxPlaybackID, _ = self.audioManager.PlaySound(self.sfxID, 0, self.stackConfig)
			}
		}
	}
}

type DebugDrawSystem struct{}

func (self *DebugDrawSystem) Draw(world *lazyecs.World, renderer *katsu2d.BatchRenderer) {
	screen := renderer.GetScreen()

	ebitenutil.DebugPrintAt(screen, "Press [Space] to start music\nPress [S] to play sfx", 20, 20)
}

// Game implements ebiten.Game interface.
type Game struct {
	engine       *katsu2d.Engine
	audioManager *katsu2d.AudioManager
}

// NewGame creates a new Game object and sets up the engine.
func NewGame() *Game {
	g := &Game{}

	// --- Engine Setup ---
	g.engine = katsu2d.NewEngine(
		katsu2d.WithWindowSize(320, 240),
		katsu2d.WithWindowTitle("Audio Example"),
		katsu2d.WithUpdateSystem(katsu2d.NewInputSystem()),
	)

	world := g.engine.World()

	// --- Audio Setup ---
	g.audioManager = g.engine.AudioManager()
	trackID, err := g.audioManager.Load("./piano.ogg")
	if err != nil {
		log.Fatalf("failed to load audio file: %v", err)
	}
	sfxID, err := g.audioManager.Load("./wood-step.mp3")
	if err != nil {
		log.Fatalf("failed to load sfx file: %v", err)
	}

	// --- Entity Setup ---
	// Create an entity that will handle input
	inputEntity := world.CreateEntity()
	inputComponent := katsu2d.NewInputComponent(keybindings)
	lazyecs.SetComponent(world, inputEntity, *inputComponent)

	// --- System Setup ---
	g.engine.AddUpdateSystem(&AudioSystem{
		audioManager: g.audioManager,
		trackID:      trackID,
		sfxID:        sfxID,
		stackConfig: &katsu2d.StackingConfig{
			Enabled:  true,
			MaxStack: 1, // Allow up to 3 instances of the same sound
		},
	})

	g.engine.AddOverlayDrawSystem(&DebugDrawSystem{})

	return g
}

func main() {
	game := NewGame()
	if err := game.engine.Run(); err != nil {
		log.Fatal(err)
	}
}
