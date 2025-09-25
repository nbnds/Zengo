package board

import (
	"image/color"
	"math/rand"
	"zenmojo/config"
)

// Board represents the game board and its state.
type Board struct {
	grid      [][]color.Color
	selectedX int
	selectedY int
}

// New creates a new, initialized game board.
func New() *Board {
	b := &Board{
		selectedX: -1,
		selectedY: -1,
		grid:      make([][]color.Color, config.GridSize),
	}

	// Initialize grid with random colors
	counts := generateColorCounts(config.Palette, config.GridSize)
	shuffledColors := make([]color.Color, 0, config.GridSize*config.GridSize)
	for i, c := range config.Palette {
		for j := 0; j < counts[i]; j++ {
			shuffledColors = append(shuffledColors, c)
		}
	}
	rand.Shuffle(len(shuffledColors), func(i, j int) {
		shuffledColors[i], shuffledColors[j] = shuffledColors[j], shuffledColors[i]
	})

	for i := range b.grid {
		b.grid[i] = make([]color.Color, config.GridSize)
		for j := range b.grid[i] {
			b.grid[i][j] = shuffledColors[i*config.GridSize+j]
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

// HandleInput processes a mouse click at the given screen coordinates.
// It returns true if a move was made (a swap occurred).
func (b *Board) HandleInput(mouseX, mouseY int) (moveMade bool) {
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
					// Swap squares
					b.grid[b.selectedY][b.selectedX], b.grid[j][i] = b.grid[j][i], b.grid[b.selectedY][b.selectedX]
					b.selectedX = -1
					b.selectedY = -1
					return true // Move was made
				}
				return false // Click was handled, but no move was made
			}
		}
	}
	return false // Click was outside the grid
}

// generateColorCounts is a helper for New to create the initial grid state.
func generateColorCounts(palette []color.Color, GridSize int) []int {
	counts := make([]int, len(palette))
	sum := 0

	for i := range palette {
		counts[i] = rand.Intn(9) + 2 // 2 to 10
	}

	for {
		sum = 0
		for _, c := range counts {
			sum += c
		}
		if sum == 100 {
			break
		}
		idx := rand.Intn(len(palette))
		if sum > 100 {
			if counts[idx] > 2 {
				counts[idx]--
			}
		} else {
			if counts[idx] < 10 {
				counts[idx]++
			}
		}
	}
	return counts
}
