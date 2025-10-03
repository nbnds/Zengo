package board

import (
	"image/color"
	"zenmojo/config"

	"github.com/hajimehoshi/ebiten/v2"
)

// Board represents the game board and its state.
type Board struct {
	grid              [][]color.Color
	selectedX         int
	selectedY         int
	IsAnimating       bool
	AnimationProgress float64
	animationDuration float64
	animatingPiece1X  int
	animatingPiece1Y  int
	animatingPiece2X  int
	animatingPiece2Y  int
}

// New creates a new, initialized game board with a random, valid layout.
func New() *Board {
	b := &Board{
		selectedX: -1,
		selectedY: -1,
	}

	totalTiles := config.GridSize * config.GridSize

	// Generate the board in three steps:
	// 1. Calculate the sizes of all color groups
	groupSizes := generateGroupSizes(totalTiles)

	// 2. Assign colors to these groups and shuffle them
	colors := assignColorsToGroups(groupSizes)

	// 3. Create the final grid from the color array
	b.grid = createColorGrid(colors)

	return b
}

// NewFromGrid creates a new board from a pre-existing grid.
func NewFromGrid(grid [][]color.Color) *Board {
	return &Board{
		grid:      grid,
		selectedX: -1,
		selectedY: -1,
	}
}

// Grid returns the current grid state for drawing.
func (b *Board) Grid() [][]color.Color {
	return b.grid
}

// Selected returns the coordinates of the selected cell.
func (b *Board) Selected() (int, int) {
	return b.selectedX, b.selectedY
}

// AnimatingPieces returns the coordinates of the two pieces being animated.
func (b *Board) AnimatingPieces() (int, int, int, int) {
	return b.animatingPiece1X, b.animatingPiece1Y, b.animatingPiece2X, b.animatingPiece2Y
}

// AnimationDuration returns the duration of the current animation in seconds.
func (b *Board) AnimationDuration() float64 {
	return b.animationDuration
}

// HandleInput processes a mouse click at the given screen coordinates.
// It returns true if a move was made (a swap occurred).
func (b *Board) HandleInput(mouseX, mouseY int) (moveMade bool) {
	if b.IsAnimating {
		return false
	}
	for i := 0; i < config.GridSize; i++ {
		for j := 0; j < config.GridSize; j++ {
			x := config.GridOriginX + i*(config.SquareSize+config.Gap)
			y := config.GridOriginY + j*(config.SquareSize+config.Gap)

			if mouseX >= x && mouseX < x+config.SquareSize && mouseY >= y && mouseY < y+config.SquareSize {
				if b.selectedX == -1 {
					// Select a square
					b.selectedX = i
					b.selectedY = j
				} else if b.selectedX == i && b.selectedY == j {
					// Deselect if clicking the same square
					b.selectedX = -1
					b.selectedY = -1
				} else {
					// Start animation
					b.IsAnimating = true
					b.AnimationProgress = 0
					b.animatingPiece1X = b.selectedX
					b.animatingPiece1Y = b.selectedY
					b.animatingPiece2X = i
					b.animatingPiece2Y = j

					// Calculate duration
					b.animationDuration = config.SwapAnimationDuration * config.StretchFactor

					b.selectedX = -1
					b.selectedY = -1
					return true // Move was initiated
				}
				return false
			}
		}
	}
	return false
}

// UpdateAnimation progresses the piece swapping animation.
// It returns true when the animation is finished.
func (b *Board) UpdateAnimation() (animationFinished bool) {
	if !b.IsAnimating {
		return false
	}

	// Update progress based on time
	if b.animationDuration > 0 {
		elapsed := 1.0 / float64(ebiten.TPS())
		b.AnimationProgress += elapsed / b.animationDuration
	} else {
		b.AnimationProgress = 1.0
	}

	if b.AnimationProgress >= 1.0 {
		b.AnimationProgress = 1.0
		// Swap pieces in the grid
		b.grid[b.animatingPiece1Y][b.animatingPiece1X], b.grid[b.animatingPiece2Y][b.animatingPiece2X] = b.grid[b.animatingPiece2Y][b.animatingPiece2X], b.grid[b.animatingPiece1Y][b.animatingPiece1X]
		b.IsAnimating = false
		return true // Animation finished
	}
	return false
}
