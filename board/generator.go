package board

import (
	"image/color"
	"math/rand"
	"zenmojo/config"
)

// Target distribution for group sizes
type sizeConstraint struct {
	size     int
	minCount int
	maxCount int
}

// generateGroupSizes calculates the sizes of color groups that will be placed on the board.
// It ensures that no single-stone groups are created and that group sizes are between
// minGroupSize and maxGroupSize, while targeting a specific distribution.
func generateGroupSizes(totalTiles int) []int {
	// Define target distribution
	targetDist := []sizeConstraint{
		{size: 10, minCount: 0, maxCount: 2}, // 0-2 large groups
		{size: 9, minCount: 1, maxCount: 3},  // 1-3 groups
		{size: 8, minCount: 1, maxCount: 4},  // 1-4 groups
		{size: 7, minCount: 1, maxCount: 3},  // 1-3 groups
		{size: 6, minCount: 2, maxCount: 4},  // 2-4 groups
		{size: 5, minCount: 2, maxCount: 4},  // 2-4 groups
		{size: 4, minCount: 2, maxCount: 5},  // 2-5 groups
		{size: 3, minCount: 2, maxCount: 4},  // 2-4 groups
		{size: 2, minCount: 1, maxCount: 4},  // 1-4 groups
	}

	var groupSizes []int
	remainingTiles := totalTiles

	minGroupSize := 2
	maxGroupSize := 10

	// First pass: ensure minimum counts for each size
	for _, sc := range targetDist {
		count := sc.minCount
		for i := 0; i < count && remainingTiles >= sc.size; i++ {
			groupSizes = append(groupSizes, sc.size)
			remainingTiles -= sc.size
		}
	}

	// Second pass: add additional groups within constraints
	for remainingTiles >= minGroupSize {
		// Try each size in random order
		indices := rand.Perm(len(targetDist))
		added := false

		for _, idx := range indices {
			sc := targetDist[idx]
			currentCount := 0
			for _, size := range groupSizes {
				if size == sc.size {
					currentCount++
				}
			}

			if currentCount < sc.maxCount && remainingTiles >= sc.size {
				groupSizes = append(groupSizes, sc.size)
				remainingTiles -= sc.size
				added = true
				break
			}
		}

		// If we couldn't add any size within constraints, add smallest possible
		if !added {
			size := minGroupSize
			for s := minGroupSize + 1; s <= maxGroupSize; s++ {
				if s <= remainingTiles {
					count := 0
					for _, gs := range groupSizes {
						if gs == s {
							count++
						}
					}
					for _, sc := range targetDist {
						if sc.size == s && count < sc.maxCount {
							size = s
							break
						}
					}
				}
			}
			groupSizes = append(groupSizes, size)
			remainingTiles -= size
		}
	}

	// Final validation: Ensure total tiles is exactly what we want
	totalSize := 0
	for _, size := range groupSizes {
		totalSize += size
	}

	// If we have too many tiles, reduce the largest groups while respecting constraints
	for totalSize > totalTiles {
		// Find largest group that can be reduced
		maxIdx := -1
		for i := 0; i < len(groupSizes); i++ {
			size := groupSizes[i]
			if size > minGroupSize {
				// Check if reducing this group still meets minimum count
				for _, sc := range targetDist {
					if sc.size == size {
						count := 0
						for _, gs := range groupSizes {
							if gs == size {
								count++
							}
						}
						if count > sc.minCount {
							if maxIdx == -1 || groupSizes[i] > groupSizes[maxIdx] {
								maxIdx = i
							}
						}
						break
					}
				}
			}
		}
		if maxIdx != -1 {
			groupSizes[maxIdx]--
			totalSize--
		} else {
			break // Can't reduce any more while respecting constraints
		}
	}

	// If we have too few tiles, add to smaller groups while respecting constraints
	for totalSize < totalTiles {
		// Find smallest group that can grow within constraints
		minIdx := -1
		for i := 0; i < len(groupSizes); i++ {
			size := groupSizes[i]
			if size < maxGroupSize {
				// Check if growing this group still meets maximum count
				for _, sc := range targetDist {
					if sc.size == size+1 {
						count := 0
						for _, gs := range groupSizes {
							if gs == size+1 {
								count++
							}
						}
						if count < sc.maxCount {
							if minIdx == -1 || groupSizes[i] < groupSizes[minIdx] {
								minIdx = i
							}
						}
						break
					}
				}
			}
		}
		if minIdx != -1 {
			groupSizes[minIdx]++
			totalSize++
		} else {
			break // Can't grow any more while respecting constraints
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
