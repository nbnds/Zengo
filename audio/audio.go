package audio

import (
	"bytes"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	// This is the target sample rate of the audio context. All sounds will be resampled to this rate before playing.
	contextSampleRate = 44100
)

// Manager holds all audio-related state.
type Manager struct {
	audioContext              *audio.Context
	moveSoundData             []byte // Raw bytes of the .wav file
	moveSoundLength           int64  // Length of the decoded PCM stream in bytes
	moveSoundNativeSampleRate int    // Sample rate of the original .wav file
}

// NewManager creates a new audio manager and loads sounds.
func NewManager(moveSoundData []byte) *Manager {
	m := &Manager{}
	// The audio context dictates the final output sample rate.
	m.audioContext = audio.NewContext(contextSampleRate)

	m.moveSoundData = moveSoundData

	// Decode the sound once (without resampling) to get its original properties.
	s, err := wav.DecodeWithoutResampling(bytes.NewReader(m.moveSoundData))
	if err != nil {
		log.Fatalf("audio: failed to decode move sound properties: %v", err)
	}

	// Store the native properties of the sound.
	m.moveSoundNativeSampleRate = s.SampleRate()
	m.moveSoundLength = s.Length()

	return m
}

// PlayMoveSound plays the sound for a piece move, adjusting its speed to the given duration.
func (m *Manager) PlayMoveSound(animationDuration float64) {
	// Decode the WAV data from memory to get a PCM stream.
	decodedStream, err := wav.DecodeWithoutResampling(bytes.NewReader(m.moveSoundData))
	if err != nil {
		log.Printf("audio: failed to decode move sound for playback: %v", err)
		return
	}

	// The stream that will be passed to the player.
	var streamToPlay io.ReadSeeker

	// We assume the audio is 16-bit stereo (4 bytes per sample frame).
	const bytesPerSampleFrame = 4
	originalDuration := float64(m.moveSoundLength) / float64(m.moveSoundNativeSampleRate*bytesPerSampleFrame)

	if animationDuration > 0 {
		// To change the speed of the sound, we can change the sample rate.
		// By telling the Resample function that the source has a different sample rate
		// than it actually does, we can trick it into producing more or fewer samples,
		// effectively changing the playback speed and pitch.
		speedRatio := originalDuration / animationDuration
		effectiveSourceSampleRate := int(float64(m.moveSoundNativeSampleRate) * speedRatio)

		// Resample the stream from its "effective" sample rate to the context's target rate.
		streamToPlay = audio.Resample(decodedStream, m.moveSoundLength, effectiveSourceSampleRate, contextSampleRate)
	} else {
		// If no animation duration is provided, just play the sound at its original speed.
		// We still need to resample it to match the audio context's sample rate if they differ.
		streamToPlay = audio.Resample(decodedStream, m.moveSoundLength, m.moveSoundNativeSampleRate, contextSampleRate)
	}

	// Create a player with the (potentially resampled) stream.
	player, err := m.audioContext.NewPlayer(streamToPlay)
	if err != nil {
		log.Printf("audio: failed to create player: %v", err)
		return
	}

	player.Play()
}
