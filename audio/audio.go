package audio

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

// Manager holds all audio-related state.
type Manager struct {
	audioContext *audio.Context
	MovePlayer   *audio.Player
}

// NewManager creates a new audio manager and loads sounds.
func NewManager() *Manager {
	var err error
	m := &Manager{}

	m.audioContext = audio.NewContext(sampleRate)

	f, err := os.Open("assets/move.wav")
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		log.Fatal(err)
	}

	m.MovePlayer, err = m.audioContext.NewPlayer(d)
	if err != nil {
		log.Fatal(err)
	}

	return m
}


// PlayMoveSound plays the sound for a piece move.
func (m *Manager) PlayMoveSound() {
	m.MovePlayer.Rewind()
	m.MovePlayer.Play()
}
