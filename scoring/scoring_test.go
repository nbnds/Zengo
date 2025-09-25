package scoring

import (
	"image/color"
	"reflect"
	"sort"
	"testing"
)

// TestFindGroups checks the core group-finding logic.
func TestFindGroups(t *testing.T) {
	// Define some colors for our test cases.
	red := color.Gray{Y: 1}
	blue := color.Gray{Y: 2}
	green := color.Gray{Y: 3}
	nilColor := color.Color(nil) // Represents an empty cell

	testCases := []struct {
		name                string
		grid                [][]color.Color
		expectedGroupSizes []int // We check the sizes of the found groups
	}{
		{
			name: "Empty grid",
			grid: [][]color.Color{
				{nilColor, nilColor},
				{nilColor, nilColor},
			},
			expectedGroupSizes: []int{},
		},
		{
			name: "Grid with single items",
			grid: [][]color.Color{
				{red, nilColor},
				{nilColor, blue},
			},
			expectedGroupSizes: []int{}, // Groups must have at least 2 items
		},
		{
			name: "Single horizontal group of 3",
			grid: [][]color.Color{
				{red, red, red},
				{nilColor, nilColor, nilColor},
			},
			expectedGroupSizes: []int{3},
		},
		{
			name: "Single vertical group of 2",
			grid: [][]color.Color{
				{blue, nilColor},
				{blue, nilColor},
			},
			expectedGroupSizes: []int{2},
		},
		{
			name: "Two separate groups",
			grid: [][]color.Color{
				{red, red, nilColor},
				{nilColor, nilColor, nilColor},
				{nilColor, blue, blue},
			},
			expectedGroupSizes: []int{2, 2},
		},
		{
			name: "Diagonally adjacent items are not a group",
			grid: [][]color.Color{
				{green, nilColor, nilColor},
				{nilColor, green, nilColor},
				{nilColor, nilColor, nilColor},
			},
			expectedGroupSizes: []int{},
		},
		{
			name: "L-shaped group of 4",
			grid: [][]color.Color{
				{red, red, red},
				{red, nilColor, nilColor},
				{nilColor, nilColor, nilColor},
			},
			expectedGroupSizes: []int{4},
		},
		{
			name: "Items of same color not touching",
			grid: [][]color.Color{
				{blue, nilColor, blue},
				{nilColor, nilColor, nilColor},
			},
			expectedGroupSizes: []int{}, // Each is a group of 1, which is ignored
		},
		{
			name: "Large group should not be split into sub-groups",
			grid: [][]color.Color{
				{red, red, red},
				{red, red, red},
				{red, red, red},
			},
			expectedGroupSizes: []int{9},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundGroups := findGroups(tc.grid)

			if len(foundGroups) != len(tc.expectedGroupSizes) {
				t.Fatalf("Expected to find %d groups, but found %d", len(tc.expectedGroupSizes), len(foundGroups))
			}

			// Check the size of each found group
			foundSizes := []int{}
			for _, g := range foundGroups {
				foundSizes = append(foundSizes, len(g.Coordinates))
			}

			// Sort both slices to make the comparison order-independent
			sort.Ints(foundSizes)
			sort.Ints(tc.expectedGroupSizes)

			if !reflect.DeepEqual(foundSizes, tc.expectedGroupSizes) {
				t.Errorf("Expected group sizes %v, but got %v", tc.expectedGroupSizes, foundSizes)
			}
		})
	}
}

func TestStandardRuleSet_Calculate(t *testing.T) {
	red := color.Gray{Y: 1}
	blue := color.Gray{Y: 2}
	nilColor := color.Color(nil)

	rules := StandardRuleSet{}

	testCases := []struct {
		name          string
		grid          [][]color.Color
		expectedScore int
	}{
		{
			name: "Correct score for horizontal line of 3",
			grid: [][]color.Color{
				{red, red, red},
				{nilColor, nilColor, nilColor},
			},
			expectedScore: 3,
		},
		{
			name: "Correct score for 2x2 square",
			grid: [][]color.Color{
				{blue, blue},
				{blue, blue},
			},
			expectedScore: 16, // 4 items * 2 width * 2 height
		},
		{
			name: "Incomplete group scores 0",
			grid: [][]color.Color{
				{red, red, nilColor},
				{red, nilColor, red}, // 4th red stone is separate
			},
			expectedScore: 0,
		},
		{
			name: "L-shaped group scores 0",
			grid: [][]color.Color{
				{blue, blue, blue},
				{blue, nilColor, nilColor},
			},
			expectedScore: 0, // Not a solid rectangle
		},
		{
			name: "Group with hole scores 0",
			grid: [][]color.Color{
				{red, red, red},
				{red, nilColor, red},
				{red, red, red},
			},
			expectedScore: 0, // Not a solid rectangle
		},
		{
			name: "Multiple valid groups are summed up",
			grid: [][]color.Color{
				{red, red, nilColor},
				{nilColor, nilColor, nilColor},
				{blue, blue, blue},
			},
			expectedScore: 5, // 2 for red group + 3 for blue group
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We test the whole CalculateScore function as it uses the rule set.
			// This is a mini-integration test for the scoring package.
			actualScore := CalculateScore(tc.grid, rules)
			if actualScore != tc.expectedScore {
				t.Errorf("Expected score %d, but got %d", tc.expectedScore, actualScore)
			}
		})
	}
}
