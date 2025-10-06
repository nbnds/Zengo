package board

import (
	"image/color"
	"zenmojo/config"
)

// BoardProperties defines the properties that any valid board must satisfy
type BoardProperties struct {
	HasNoSingleStones  bool
	HasValidGroupSizes bool
	UsesValidColors    bool
	IsFull             bool
}

// Helper function to compare colors (since color.Color doesn't implement equality)
func colorEqual(c1, c2 color.Color) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// checkBoardProperties verifies all required properties of a game board
func checkBoardProperties(grid [][]color.Color) BoardProperties {
	// Initialize properties as true (assume all checks will pass)
	props := BoardProperties{
		HasNoSingleStones:  true,
		HasValidGroupSizes: true,
		UsesValidColors:    true,
		IsFull:             true,
	}

	// Count stones by color
	colorCounts := make(map[color.Color]int)
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			c := grid[i][j]
			if c == nil {
				props.IsFull = false
				continue
			}
			colorCounts[c]++

			// Check if color is from valid palette
			validColor := false
			for _, paletteColor := range config.Palette {
				if colorEqual(c, paletteColor) {
					validColor = true
					break
				}
			}
			if !validColor {
				props.UsesValidColors = false
			}
		}
	}

	// Check group sizes
	for _, count := range colorCounts {
		if count == 1 {
			props.HasNoSingleStones = false
		}
		if count < 2 || count > 10 {
			props.HasValidGroupSizes = false
		}
	}

	return props
}
