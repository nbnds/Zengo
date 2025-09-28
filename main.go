package main

import (
	_ "embed"
	"log"
	"zenmojo/audio"
	"zenmojo/config"
	"zenmojo/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/move.wav
var moveSoundFile []byte

func main() {
	// Set window properties from the config package
	ebiten.SetWindowIcon(config.Icons)
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Zengo")

	// Create a new audio manager
	audioManager := audio.NewManager(moveSoundFile)

	// Create and run the game
	game := game.NewGame(audioManager)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
