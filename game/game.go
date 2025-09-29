package game

import (
	"image/color"
	"log"
	"time"
	"zenmojo/audio"
	"zenmojo/board"
	"zenmojo/config"
	"zenmojo/scoring"
	"zenmojo/sharing"
	"zenmojo/view"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.design/x/clipboard"
)

// Game holds the main game state.
type Game struct {
	board            *board.Board
	audioManager     *audio.Manager
	score            int
	maxScore         int
	moveCount        int
	scoreHistory     []int
	colorCounts      map[color.Color]int
	shareCode        string
	isCustomBoard    bool
	copyFeedback     string
	copyFeedbackTime time.Time
}

// NewGame initializes a new game.
func NewGame(audioManager *audio.Manager) *Game {
	g := &Game{
		audioManager: audioManager,
	}
	g.startNewGame(nil) // Start with a random board
	return g
}

// startNewGame resets the game state with a new board.
// If a grid is provided, it uses that; otherwise, it creates a random one.
func (g *Game) startNewGame(grid [][]color.Color) {
	if grid == nil {
		g.board = board.New()
		g.isCustomBoard = false
	} else {
		g.board = board.NewFromGrid(grid)
		g.isCustomBoard = true
	}

	// Update the window icon to match a tile from the new board.
	if g.board.Grid()[0][0] != nil {
		ebiten.SetWindowIcon(config.CreateTileIcons(g.board.Grid()[0][0]))
	}

	g.score = scoring.CalculateScore(g.board.Grid(), scoring.StandardRuleSet{})
	g.maxScore = scoring.CalculateMaxPossibleScore(g.board.Grid())
	g.moveCount = 0
	g.scoreHistory = []int{g.score}
	g.colorCounts = scoring.CountColors(g.board.Grid())

	// Generate the share code for this board
	code, err := sharing.Encode(g.board.Grid())
	if err != nil {
		log.Printf("Error generating share code: %v", err)
		g.shareCode = "Error"
	} else {
		g.shareCode = code
	}
}

// Update proceeds the game state.
func (g *Game) Update() error {
	// Check for pasted share code
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		pastedText := string(clipboard.Read(clipboard.FmtText))
		if grid, err := sharing.Decode(pastedText); err == nil {
			g.startNewGame(grid)
			return nil // Restarted game, skip rest of update
		}
	}

	// Handle board animation
	if g.board.IsAnimating {
		if g.board.UpdateAnimation() {
			// Animation finished, recalculate score
			g.score = scoring.CalculateScore(g.board.Grid(), scoring.StandardRuleSet{})
			g.scoreHistory = append(g.scoreHistory, g.score)
		}
		return nil
	}

	// Handle mouse input for piece selection/swapping
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Check if the share code was clicked
		if view.IsShareCodeClicked(x, y, g.shareCode) {
			clipboard.Write(clipboard.FmtText, []byte(g.shareCode))
			g.copyFeedback = "Copied!"
			g.copyFeedbackTime = time.Now()
		} else if g.board.HandleInput(x, y) {
			// A move was made
			g.moveCount++
			g.audioManager.PlayMoveSound(g.board.AnimationDuration())
		}
	}

	// Reset copy feedback message after a delay
	if g.copyFeedback != "" && time.Since(g.copyFeedbackTime).Seconds() > 1.5 {
		g.copyFeedback = ""
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	mouseX, mouseY := ebiten.CursorPosition()
	view.Draw(screen, g.board, g.score, g.maxScore, g.moveCount, g.scoreHistory, g.colorCounts, mouseX, mouseY)

	// Draw the new sharing UI elements
	view.DrawSharingUI(screen, g.shareCode, g.copyFeedback, g.isCustomBoard)
}

// Layout is called when the window is resized.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}
