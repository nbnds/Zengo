package main

import (
	"testing"
)

func TestGenerateColorCounts(t *testing.T) {
	counts := generateColorCounts(palette, gridSize)

	sum := 0
	for _, count := range counts {
		sum += count
	}

	if sum != 100 {
		t.Errorf("Sum of counts should be 100, but got %d", sum)
	}

	for i, count := range counts {
		if count < 2 || count > 10 {
			t.Errorf("Count for color %d should be between 2 and 10, but got %d", i, count)
		}
	}

	notUniform := false
	for _, count := range counts {
		if count != 10 {
			notUniform = true
			break
		}
	}

	if !notUniform {
		t.Errorf("All counts are 10, the distribution is uniform")
	}
}
