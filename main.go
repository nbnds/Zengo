package main

import (
	_ "embed"
	"errors"
	"log"
	"zenmojo/audio"
	"zenmojo/config"
	"zenmojo/game"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.design/x/clipboard"
)

//go:embed assets/move.wav
var moveSoundFile []byte

func main() {
	// Initialize the cross-platform clipboard package.
	err := clipboard.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Set window properties from the config package
	ebiten.SetWindowIcon(config.Icons)
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Zengo")

	// Create a new audio manager
	audioManager := audio.NewManager(moveSoundFile)

	// Create and run the game
	game := game.NewGame(audioManager)
	if err := ebiten.RunGame(game); err != nil {
		// ebiten.Termination is a sentinel error indicating a clean exit.
		if !errors.Is(err, ebiten.Termination) {
			log.Fatal(err)
		}
	}
}
