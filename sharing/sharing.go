package sharing

import (
	"errors"
	"image/color"
	"strings"
	"zenmojo/config"
)

// Define the character set for encoding. Using a URL-safe set is good practice.
const encodingChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

var (
	colorToChar   map[string]rune
	charToColor   map[rune]color.Color
	isInitialized = false
)

// initialize prepares the mapping tables. This is done once.
func initialize() {
	if isInitialized || len(config.Palette) > len(encodingChars) {
		return
	}
	colorToChar = make(map[string]rune)
	charToColor = make(map[rune]color.Color)

	for i, c := range config.Palette {
		// Use a string representation of the color as a map key
		r, g, b, a := c.RGBA()
		key := string(r) + string(g) + string(b) + string(a)

		char := rune(encodingChars[i])
		colorToChar[key] = char
		charToColor[char] = c
	}
	isInitialized = true
}

// Encode takes a board grid and converts it into a shareable string code.
func Encode(grid [][]color.Color) (string, error) {
	initialize()
	if !isInitialized {
		return "", errors.New("sharing: palette size exceeds encoding character set")
	}

	var sb strings.Builder
	sb.Grow(config.GridSize * config.GridSize)

	for r := 0; r < config.GridSize; r++ {
		for c := 0; c < config.GridSize; c++ {
			cellColor := grid[r][c]
			r, g, b, a := cellColor.RGBA()
			key := string(r) + string(g) + string(b) + string(a)

			char, ok := colorToChar[key]
			if !ok {
				return "", errors.New("sharing: color not found in palette")
			}
			sb.WriteRune(char)
		}
	}
	return sb.String(), nil
}

// Decode takes a shareable code and converts it back into a board grid.
func Decode(code string) ([][]color.Color, error) {
	initialize()
	if len(code) != config.GridSize*config.GridSize {
		return nil, errors.New("sharing: invalid code length")
	}

	grid := make([][]color.Color, config.GridSize)
	for r := 0; r < config.GridSize; r++ {
		grid[r] = make([]color.Color, config.GridSize)
		for c := 0; c < config.GridSize; c++ {
			char := rune(code[r*config.GridSize+c])
			grid[r][c] = charToColor[char]
		}
	}
	return grid, nil
}
