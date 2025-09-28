package scoring

import "image/color"

// ScoringRule defines the contract for a set of scoring rules.
// This allows for different scoring strategies (e.g., for different game modes).
type ScoringRule interface {
	Calculate(group Group, grid [][]color.Color) int
}

// StandardRuleSet implements the default scoring logic where only
// complete, solid rectangles score points.
type StandardRuleSet struct{}

// Calculate applies the standard scoring rules to a single group.
func (s StandardRuleSet) Calculate(group Group, grid [][]color.Color) int {
	// This is the logic from the old calculateScoreForGroup function.
	// We will add the new "solid rectangle" rule here in a later step.
	totalItemsOnBoard := totalColorItems(group.Color, grid)
	if len(group.Coordinates) != totalItemsOnBoard {
		return 0 // Only score groups that contain all items of that color.
	}

	if len(group.Coordinates) < 2 {
		return 0
	}

	width := group.MaxC - group.MinC + 1
	height := group.MaxR - group.MinR + 1
	numItems := len(group.Coordinates)

	// Line shape: score is the number of items.
	if width == 1 || height == 1 {
		return numItems
	}

	// Check if the group is a solid rectangle. If not, it scores 0.
	if numItems != width*height {
		return 0
	}

	// For solid rectangles (including squares), score is items * width * height.
	return numItems * width * height
}

// Coordinate represents a position on the grid.
type Coordinate struct {
	R, C int
}

// Group represents a contiguous group of same-colored items.
type Group struct {
	Color       color.Color
	Coordinates []Coordinate
	MinR, MaxR  int // Bounding box for shape detection
	MinC, MaxC  int // Bounding box for shape detection
}

// CalculateScore analyzes the entire grid using a given rule set and returns the total score.
func CalculateScore(grid [][]color.Color, rule ScoringRule) int {
	groups := findGroups(grid)
	totalScore := 0
	for _, group := range groups {
		totalScore += rule.Calculate(group, grid)
	}
	return totalScore
}

func findGroups(grid [][]color.Color) []Group {
	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	var groups []Group

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if visited[r][c] || grid[r][c] == nil {
				continue
			}

			currentGroup := Group{
				Color:       grid[r][c],
				Coordinates: []Coordinate{},
				MinR:        r,
				MaxR:        r,
				MinC:        c,
				MaxC:        c,
			}

			dfs(r, c, grid[r][c], &currentGroup, visited, grid)

			if len(currentGroup.Coordinates) >= 2 {
				groups = append(groups, currentGroup)
			}
		}
	}
	return groups
}

func dfs(r, c int, targetColor color.Color, currentGroup *Group, visited [][]bool, grid [][]color.Color) {
	rows, cols := len(grid), len(grid[0])
	if r < 0 || r >= rows || c < 0 || c >= cols || visited[r][c] || !colorsEqual(grid[r][c], targetColor) {
		return
	}

	visited[r][c] = true
	currentGroup.Coordinates = append(currentGroup.Coordinates, Coordinate{R: r, C: c})

	if r < currentGroup.MinR {
		currentGroup.MinR = r
	}
	if r > currentGroup.MaxR {
		currentGroup.MaxR = r
	}
	if c < currentGroup.MinC {
		currentGroup.MinC = c
	}
	if c > currentGroup.MaxC {
		currentGroup.MaxC = c
	}

	directions := []struct{ dr, dc int }{
		{-1, 0}, // Up
		{1, 0},  // Down
		{0, -1}, // Left
		{0, 1},  // Right
	}
	for _, dir := range directions {
		dfs(r+dir.dr, c+dir.dc, targetColor, currentGroup, visited, grid)
	}
}

func colorsEqual(c1, c2 color.Color) bool {
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

func totalColorItems(color color.Color, grid [][]color.Color) int {
	count := 0
	for r := range grid {
		for c := range grid[r] {
			if colorsEqual(grid[r][c], color) {
				count++
			}
		}
	}
	return count
}

// CalculateMaxPossibleScore determines the theoretical maximum score for a given board layout.
// It does this by counting the items of each color and calculating the score for the
// most optimal shape (the most "square-like" rectangle) that can be formed with that number of items.
func CalculateMaxPossibleScore(grid [][]color.Color) int {
	colorCounts := make(map[color.Color]int)
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] != nil {
				colorCounts[grid[r][c]]++
			}
		}
	}

	totalMaxScore := 0
	for _, numItems := range colorCounts {
		if numItems < 2 {
			continue
		}

		bestScoreForColor := 0
		// Find all factors to determine possible rectangle shapes
		for w := 1; w*w <= numItems; w++ {
			if numItems%w == 0 {
				h := numItems / w
				var score int
				// Apply the correct scoring rule based on the shape.
				if w == 1 || h == 1 {
					// For a line shape, the score is just the number of items.
					score = numItems
				} else {
					// For a solid rectangle, the score is items * width * height.
					score = numItems * w * h
				}

				if score > bestScoreForColor {
					bestScoreForColor = score
				}
			}
		}
		totalMaxScore += bestScoreForColor
	}
	return totalMaxScore
}

// CountColors counts the number of tiles for each color on the grid.
func CountColors(grid [][]color.Color) map[color.Color]int {
	colorCounts := make(map[color.Color]int)
	for r := range grid {
		for c := range grid[r] {
			if grid[r][c] != nil {
				colorCounts[grid[r][c]]++
			}
		}
	}
	return colorCounts
}
