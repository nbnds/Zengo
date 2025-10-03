package board

import (
	"image/color"
	"testing"
	"zenmojo/config"
	"zenmojo/sharing"
)

func TestNew(t *testing.T) {
	// Run board generation multiple times to ensure consistency
	numTrials := 100

	for trial := 0; trial < numTrials; trial++ {
		board := New()
		grid := board.Grid()

		t.Run("basic board properties", func(t *testing.T) {
			// Check board dimensions
			if len(grid) != config.GridSize {
				t.Errorf("Expected grid height %d, got %d", config.GridSize, len(grid))
			}
			for row := 0; row < config.GridSize; row++ {
				if len(grid[row]) != config.GridSize {
					t.Errorf("Row %d: expected width %d, got %d", row, config.GridSize, len(grid[row]))
				}
			}

			// Check for nil cells
			for i := 0; i < config.GridSize; i++ {
				for j := 0; j < config.GridSize; j++ {
					if grid[i][j] == nil {
						t.Errorf("Found nil cell at position (%d, %d)", i, j)
					}
				}
			}
		})

		t.Run("group size constraints", func(t *testing.T) {
			// Create a map to count stones of each color
			colorCounts := make(map[color.Color]int)
			for i := 0; i < config.GridSize; i++ {
				for j := 0; j < config.GridSize; j++ {
					c := grid[i][j]
					colorCounts[c]++
				}
			}

			// Check that no color has just one stone
			for c, count := range colorCounts {
				if count == 1 {
					// Get share code for reproduction
					shareCode, _ := sharing.Encode(grid)
					t.Errorf("Found color with single stone (count=1): %v\nBoard share code: %s", c, shareCode)
				}
			}

			// Check that no group exceeds maxGroupSize
			maxGroupSize := 10
			for _, count := range colorCounts {
				if count > maxGroupSize {
					t.Errorf("Found group larger than maximum size. Size: %d, Max allowed: %d", count, maxGroupSize)
				}
			}
		})

		t.Run("color distribution", func(t *testing.T) {
			// Track which colors from the palette are used
			usedColors := make(map[color.Color]bool)
			for i := 0; i < config.GridSize; i++ {
				for j := 0; j < config.GridSize; j++ {
					usedColors[grid[i][j]] = true
				}
			}

			// Ensure all used colors are from the palette
			for c := range usedColors {
				found := false
				for _, paletteColor := range config.Palette {
					if colorEqual(c, paletteColor) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Found color not in palette: %v", c)
				}
			}
		})
	}
}

// Test that we can consistently decode and encode board states
func TestBoardEncodeDecode(t *testing.T) {
	// Create a new board
	originalBoard := New()
	originalGrid := originalBoard.Grid()

	// Create a new board with the same grid
	copiedBoard := NewFromGrid(originalGrid)
	copiedGrid := copiedBoard.Grid()

	// Compare the grids
	for i := 0; i < config.GridSize; i++ {
		for j := 0; j < config.GridSize; j++ {
			if !colorEqual(originalGrid[i][j], copiedGrid[i][j]) {
				t.Errorf("Grid mismatch at position (%d, %d)", i, j)
			}
		}
	}
}
