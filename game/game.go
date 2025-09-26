package game

import (
	"zenmojo/audio"
	"zenmojo/board"
	"zenmojo/config"
	"zenmojo/scoring"
	"zenmojo/view"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game holds all game state and logic.
type Game struct {
	board        *board.Board
	mouseX       int
	mouseY       int
	moveCount    int
	score        int
	audioManager *audio.Manager
}

// NewGame creates a new, initialized Game object.
func NewGame(audioManager *audio.Manager) *Game {
	b := board.New()
	g := &Game{
		board:        b,
		moveCount:    0,
		score:        scoring.CalculateScore(b.Grid(), scoring.StandardRuleSet{}), // Calculate initial score
		audioManager: audioManager,
	}
	return g
}

// Update handles the game logic for each frame.
func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	if animationFinished := g.board.UpdateAnimation(); animationFinished {
		g.score = scoring.CalculateScore(g.board.Grid(), scoring.StandardRuleSet{})
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if moveInitiated := g.board.HandleInput(g.mouseX, g.mouseY); moveInitiated {
			g.moveCount++
			g.audioManager.PlayMoveSound()
		}
	}
	return nil
}

// Draw renders the game screen by delegating to the view package.
func (g *Game) Draw(screen *ebiten.Image) {
	view.Draw(screen, g.board, g.score, g.moveCount, g.mouseX, g.mouseY)
}

// Layout returns the configured screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
