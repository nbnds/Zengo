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

	// Verteile die Steine in Gruppen
	for remainingTiles > 0 {
		// Berechne die minimale Anzahl weiterer Gruppen, die wir noch brauchen
		minRemainingGroups := (remainingTiles + maxGroupSize - 1) / maxGroupSize
		if minRemainingGroups < remainingTiles/maxGroupSize {
			minRemainingGroups++
		}

		// Berechne die Größe für diese Gruppe
		size := maxGroupSize
		if remainingTiles < maxGroupSize {
			size = remainingTiles
		}

		// Wenn die verbleibenden Steine zu wenige für eine neue Gruppe sind,
		// verteile sie auf existierende Gruppen
		if remainingTiles-size < minGroupSize && remainingTiles-size > 0 {
			extraStones := remainingTiles - size

			// Versuche die extra Steine auf vorhandene Gruppen zu verteilen
			for i := len(groupSizes) - 1; i >= 0 && extraStones > 0; i-- {
				spaceLeft := maxGroupSize - groupSizes[i]
				if spaceLeft > 0 {
					add := extraStones
					if add > spaceLeft {
						add = spaceLeft
					}
					groupSizes[i] += add
					extraStones -= add
				}
			}

			// Wenn noch Steine übrig sind, füge sie zur aktuellen Gruppe hinzu
			if extraStones > 0 {
				size += extraStones
			}
		}

		// Füge die neue Gruppe hinzu
		if size >= minGroupSize {
			groupSizes = append(groupSizes, size)
			remainingTiles -= size
		} else {
			// Wenn die Gruppe zu klein wäre, verteile die Steine auf vorhandene Gruppen
			stonesLeft := size
			for i := len(groupSizes) - 1; i >= 0 && stonesLeft > 0; i-- {
				spaceLeft := maxGroupSize - groupSizes[i]
				if spaceLeft > 0 {
					add := stonesLeft
					if add > spaceLeft {
						add = spaceLeft
					}
					groupSizes[i] += add
					stonesLeft -= add
				}
			}
			// Wenn wir die Steine nicht verteilen konnten, müssen wir eine neue Gruppe machen
			if stonesLeft > 0 {
				// Hole mehr Steine von der letzten Gruppe
				lastGroupIdx := len(groupSizes) - 1
				if lastGroupIdx >= 0 {
					// Nimm genug Steine von der letzten Gruppe, um eine gültige neue Gruppe zu bilden
					needed := minGroupSize - stonesLeft
					if groupSizes[lastGroupIdx] > needed+minGroupSize {
						groupSizes[lastGroupIdx] -= needed
						stonesLeft += needed
						groupSizes = append(groupSizes, stonesLeft)
					}
				}
			}
			remainingTiles -= size
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
