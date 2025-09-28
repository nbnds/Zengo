package board

import (
	"image/color"
	"math/rand"
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
		grid:      make([][]color.Color, config.GridSize),
	}

	// --- New distribution logic for varied group sizes ---
	totalTiles := config.GridSize * config.GridSize
	minGroupSize := 2
	maxGroupSize := 10

	colors := make([]color.Color, 0, totalTiles)
	remainingTiles := totalTiles

	// Create a temporary, mutable copy of the palette to ensure each color is used at most once.
	availableColors := make([]color.Color, len(config.Palette))
	copy(availableColors, config.Palette)
	rand.Shuffle(len(availableColors), func(i, j int) {
		availableColors[i], availableColors[j] = availableColors[j], availableColors[i]
	})

	// Create groups of random sizes
	for remainingTiles > 0 {
		// Determine a random size for the next group
		size := rand.Intn(maxGroupSize-minGroupSize+1) + minGroupSize

		// If we've run out of unique colors, stop creating new groups.
		// The last group will take all remaining tiles.
		// Ensure the last group doesn't leave a single tile
		if remainingTiles-size == 1 {
			size = remainingTiles / 2 // Split the remainder
		}
		if remainingTiles-size < 0 {
			size = remainingTiles // Last group takes all remaining tiles
		}
		if len(availableColors) == 0 {
			// If we are out of colors, distribute the rest among the last color,
			// but respect maxGroupSize.
			if remainingTiles > maxGroupSize {
				size = maxGroupSize
			} else {
				size = remainingTiles
			}
		}

		// Assign a unique color to this group and remove it from the available pool.
		// If no unique colors are left, reuse the last available color for the final group.
		groupColor := availableColors[0]
		if len(availableColors) > 1 {
			availableColors = availableColors[1:]
		}
		for i := 0; i < size; i++ {
			colors = append(colors, groupColor)
		}

		remainingTiles -= size
	}

	// Shuffle the colors to create a random board.
	// This ensures that the generated groups are scattered across the grid.
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	// Populate the grid with the shuffled colors.
	for i := 0; i < config.GridSize; i++ {
		b.grid[i] = make([]color.Color, config.GridSize)
		for j := 0; j < config.GridSize; j++ {
			b.grid[i][j] = colors[i*config.GridSize+j]
		}
	}

	return b
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
