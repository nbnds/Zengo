package main

import (
	"log"
	"zenmojo/config"
	"zenmojo/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Set window properties from the config package
	ebiten.SetWindowIcon(config.Icons)
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Zesty Zen")

	// Create and run the game
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}