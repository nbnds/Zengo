package game

import (
	"math/rand"
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
	maxScore     int
	scoreHistory []int
	audioManager *audio.Manager
}

// NewGame creates a new, initialized Game object.
func NewGame(audioManager *audio.Manager) *Game {
	b := board.New()
	g := &Game{
		board:     b,
		moveCount: 0,
		score:     scoring.CalculateScore(b.Grid(), scoring.StandardRuleSet{}), // Calculate initial score
		maxScore:  scoring.CalculateMaxPossibleScore(b.Grid()),
		// Initialize score history with the score at move 0
		scoreHistory: []int{scoring.CalculateScore(b.Grid(), scoring.StandardRuleSet{})},
		audioManager: audioManager,
	}

	// --- Set a dynamic window icon based on a random tile from the new board ---
	randomRow := rand.Intn(config.GridSize)
	randomCol := rand.Intn(config.GridSize)
	iconColor := b.Grid()[randomRow][randomCol]
	ebiten.SetWindowIcon(config.CreateTileIcons(iconColor))

	return g
}

// Update handles the game logic for each frame.
func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	if animationFinished := g.board.UpdateAnimation(); animationFinished {
		newScore := scoring.CalculateScore(g.board.Grid(), scoring.StandardRuleSet{})
		g.score = newScore
		g.scoreHistory = append(g.scoreHistory, newScore)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if moveInitiated := g.board.HandleInput(g.mouseX, g.mouseY); moveInitiated {
			g.moveCount++
			duration := g.board.AnimationDuration() // Add a small buffer to ensure sound syncs with animation end
			g.audioManager.PlayMoveSound(duration)
		}
	}
	return nil
}

// Draw renders the game screen by delegating to the view package.
func (g *Game) Draw(screen *ebiten.Image) {
	colorCounts := scoring.CountColors(g.board.Grid())
	view.Draw(screen, g.board, g.score, g.maxScore, g.moveCount, g.scoreHistory, colorCounts, g.mouseX, g.mouseY)
}

// Layout returns the configured screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
