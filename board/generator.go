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

	// First pass: Create groups with size distribution, aiming for average group size
	for remainingTiles >= minGroupSize {
		// Calculate ideal average group size remaining
		avgSize := float64(remainingTiles) / float64((remainingTiles+maxGroupSize-1)/maxGroupSize)
		targetSize := int(avgSize)
		if targetSize < minGroupSize {
			targetSize = minGroupSize
		}
		if targetSize > maxGroupSize {
			targetSize = maxGroupSize
		}

		// Add some randomization around the target size
		variance := (maxGroupSize - minGroupSize) / 4
		if variance < 1 {
			variance = 1
		}

		size := targetSize + rand.Intn(variance*2+1) - variance
		if size < minGroupSize {
			size = minGroupSize
		}
		if size > maxGroupSize {
			size = maxGroupSize
		}

		// Ensure we leave enough tiles for at least one more group if tiles remain
		tilesAfterGroup := remainingTiles - size
		if tilesAfterGroup > 0 && tilesAfterGroup < minGroupSize {
			// If we can't make another group, add remaining tiles to current group
			if size+tilesAfterGroup <= maxGroupSize {
				size += tilesAfterGroup
			} else {
				// Otherwise make this group smaller to allow another valid group
				size = remainingTiles / 2
			}
		}

		groupSizes = append(groupSizes, size)
		remainingTiles -= size
	}

	// Second pass: If we have leftover tiles, distribute them evenly
	if remainingTiles > 0 {
		// Sort group sizes ascending for even distribution
		indices := make([]int, len(groupSizes))
		for i := range indices {
			indices[i] = i
		}
		for i := 0; i < len(indices)-1; i++ {
			for j := i + 1; j < len(indices); j++ {
				if groupSizes[indices[i]] > groupSizes[indices[j]] {
					indices[i], indices[j] = indices[j], indices[i]
				}
			}
		}

		// Distribute remaining tiles to smaller groups first
		for i := 0; remainingTiles > 0 && i < len(groupSizes); i++ {
			idx := indices[i]
			if groupSizes[idx] < maxGroupSize {
				add := 1 // Add only one tile at a time for more even distribution
				if add > remainingTiles {
					add = remainingTiles
				}
				if groupSizes[idx]+add <= maxGroupSize {
					groupSizes[idx] += add
					remainingTiles -= add
				}
			}
		}
	}

	// Final validation: Ensure total tiles is exactly what we want
	totalSize := 0
	for _, size := range groupSizes {
		totalSize += size
	}

	// If we have too many tiles, reduce the largest groups
	for totalSize > totalTiles {
		// Find largest group
		maxIdx := 0
		for i := 1; i < len(groupSizes); i++ {
			if groupSizes[i] > groupSizes[maxIdx] {
				maxIdx = i
			}
		}
		if groupSizes[maxIdx] > minGroupSize {
			groupSizes[maxIdx]--
			totalSize--
		}
	}

	// If we have too few tiles, add to smaller groups
	for totalSize < totalTiles {
		// Find smallest group that can grow
		minIdx := -1
		for i := 0; i < len(groupSizes); i++ {
			if groupSizes[i] < maxGroupSize {
				if minIdx == -1 || groupSizes[i] < groupSizes[minIdx] {
					minIdx = i
				}
			}
		}
		if minIdx != -1 {
			groupSizes[minIdx]++
			totalSize++
		}
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
