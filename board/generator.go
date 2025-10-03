package board

import (
	"image/color"
	"math/rand"
	"zenmojo/config"
)

// generateGroupSizes calculates the sizes of color groups that will be placed on the board.
// It ensures that no single-stone groups are created and that group sizes are between
// minGroupSize and maxGroupSize.
func generateGroupSizes(totalTiles int) []int {
	minGroupSize := 2
	maxGroupSize := 10
	remainingTiles := totalTiles
	var groupSizes []int

	for remainingTiles > 0 {
		// Determine size for the next group
		size := rand.Intn(maxGroupSize-minGroupSize+1) + minGroupSize

		// Prevent single tile remainders
		if remainingTiles-size == 1 {
			size = remainingTiles / 2 // Split the remainder
		}
		if remainingTiles-size < 0 {
			size = remainingTiles // Last group takes all remaining tiles
		}

		groupSizes = append(groupSizes, size)
		remainingTiles -= size
	}

	return groupSizes
}

// assignColorsToGroups takes the group sizes and assigns colors to each group,
// ensuring that each color is used at most once until we run out of colors.
func assignColorsToGroups(groupSizes []int) []color.Color {
	var colors []color.Color

	// Create a shuffled copy of the palette
	availableColors := make([]color.Color, len(config.Palette))
	copy(availableColors, config.Palette)
	rand.Shuffle(len(availableColors), func(i, j int) {
		availableColors[i], availableColors[j] = availableColors[j], availableColors[i]
	})

	colorIndex := 0
	for _, size := range groupSizes {
		// If we're out of colors, reuse the last color
		if colorIndex >= len(availableColors) {
			colorIndex = len(availableColors) - 1
		}

		// Create the group with the current color
		for i := 0; i < size; i++ {
			colors = append(colors, availableColors[colorIndex])
		}
		colorIndex++
	}

	// Shuffle the final color array to distribute groups randomly
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	return colors
}

// createColorGrid converts a flat array of colors into a 2D grid.
func createColorGrid(colors []color.Color) [][]color.Color {
	grid := make([][]color.Color, config.GridSize)
	for i := 0; i < config.GridSize; i++ {
		grid[i] = make([]color.Color, config.GridSize)
		for j := 0; j < config.GridSize; j++ {
			grid[i][j] = colors[i*config.GridSize+j]
		}
	}
	return grid
}
