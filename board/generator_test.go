package board

import (
	"image/color"
	"testing"
	"zenmojo/config"
)

func TestGenerateGroupSizes(t *testing.T) {
	totalTiles := config.GridSize * config.GridSize
	numTrials := 100

	for trial := 0; trial < numTrials; trial++ {
		sizes := generateGroupSizes(totalTiles)

		// Test 1: Die Summe aller Gruppengrößen muss der Gesamtzahl der Felder entsprechen
		sum := 0
		for _, size := range sizes {
			sum += size
		}
		if sum != totalTiles {
			t.Errorf("Trial %d: Sum of group sizes %d does not match total tiles %d", trial, sum, totalTiles)
		}

		// Test 2: Keine Gruppe darf kleiner als 2 oder größer als 10 sein
		for i, size := range sizes {
			if size < 2 {
				t.Errorf("Trial %d: Group %d has invalid size %d (smaller than 2)", trial, i, size)
			}
			if size > 10 {
				t.Errorf("Trial %d: Group %d has invalid size %d (larger than 10)", trial, i, size)
			}
		}
	}
}

func TestAssignColorsToGroups(t *testing.T) {
	groupSizes := []int{4, 6, 8} // Beispiel-Gruppengrößen
	totalTiles := 18             // Summe der Gruppengrößen

	colors := assignColorsToGroups(groupSizes)

	// Test 1: Die Länge der Farbliste muss der Summe der Gruppengrößen entsprechen
	if len(colors) != totalTiles {
		t.Errorf("Expected %d colors, got %d", totalTiles, len(colors))
	}

	// Test 2: Alle Farben müssen aus der Palette stammen
	for i, c := range colors {
		found := false
		for _, paletteColor := range config.Palette {
			if colorEqual(c, paletteColor) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Color at index %d is not from palette", i)
		}
	}
}

func TestCreateColorGrid(t *testing.T) {
	totalTiles := config.GridSize * config.GridSize
	colors := make([]color.Color, totalTiles)
	// Fülle das Array mit einer Testfarbe
	testColor := config.Palette[0]
	for i := range colors {
		colors[i] = testColor
	}

	grid := createColorGrid(colors)

	// Test 1: Überprüfe die Grid-Dimensionen
	if len(grid) != config.GridSize {
		t.Errorf("Expected grid height %d, got %d", config.GridSize, len(grid))
	}
	for i, row := range grid {
		if len(row) != config.GridSize {
			t.Errorf("Row %d: expected width %d, got %d", i, config.GridSize, len(row))
		}
	}

	// Test 2: Überprüfe, ob alle Farben korrekt übertragen wurden
	for i := 0; i < config.GridSize; i++ {
		for j := 0; j < config.GridSize; j++ {
			if !colorEqual(grid[i][j], colors[i*config.GridSize+j]) {
				t.Errorf("Color mismatch at position (%d,%d)", i, j)
			}
		}
	}
}
