package view

import (
	"fmt"
	"image/color"
	"sort"
	"zenmojo/board"
	"zenmojo/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

// Draw renders the entire game screen.
//
//go:noinline
func Draw(screen *ebiten.Image, b *board.Board, score, maxScore, moveCount int, scoreHistory []int, colorCounts map[color.Color]int, mouseX, mouseY int) {
	drawBackground(screen)
	drawBoard(screen, b, mouseX, mouseY)
	drawUI(screen, score, maxScore, moveCount, scoreHistory)
	drawStoneDistribution(screen, colorCounts)
}

//go:noinline
func drawBackground(screen *ebiten.Image) {
	screen.Fill(config.BackgroundColor)

	// Draw hatching pattern
	patternSize := config.HatchingPattern.Bounds().Dx()
	for i := 0; i < config.ScreenWidth; i += patternSize {
		for j := 0; j < config.ScreenHeight; j += patternSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i), float64(j))
			screen.DrawImage(config.HatchingPattern, op)
		}
	}
}

//go:noinline
func drawBoard(screen *ebiten.Image, b *board.Board, mouseX int, mouseY int) {
	if b.IsAnimating {
		// Animation logic
		p1x, p1y, p2x, p2y := b.AnimatingPieces()
		progress := b.AnimationProgress

		// Calculate interpolated positions in terms of screen coordinates
		p1StartX := float64(config.GridOriginX + p1x*(config.SquareSize+config.Gap))
		p1StartY := float64(config.GridOriginY + p1y*(config.SquareSize+config.Gap))
		p1EndX := float64(config.GridOriginX + p2x*(config.SquareSize+config.Gap))
		p1EndY := float64(config.GridOriginY + p2y*(config.SquareSize+config.Gap))

		p1CurrentX := p1StartX + (p1EndX-p1StartX)*progress
		p1CurrentY := p1StartY + (p1EndY-p1StartY)*progress

		p2CurrentX := p1EndX + (p1StartX-p1EndX)*progress
		p2CurrentY := p1EndY + (p1StartY-p1EndY)*progress

		// Draw all non-animating pieces
		for i := 0; i < config.GridSize; i++ {
			for j := 0; j < config.GridSize; j++ {
				if (i == p1x && j == p1y) || (i == p2x && j == p2y) {
					continue // Skip animating pieces, they will be drawn on top
				}
				drawPiece(screen, b, i, j, -1, -1) // Pass -1 for mouse to avoid hover effects
			}
		}

		// Draw the two animating pieces at their interpolated positions
		color1 := b.Grid()[p1y][p1x]
		color2 := b.Grid()[p2y][p2x]
		drawPieceAt(screen, color1, p1CurrentX, p1CurrentY, false)
		drawPieceAt(screen, color2, p2CurrentX, p2CurrentY, false)

	} else {
		// Original drawing logic if not animating
		for i := 0; i < config.GridSize; i++ {
			for j := 0; j < config.GridSize; j++ {
				drawPiece(screen, b, i, j, mouseX, mouseY)
			}
		}
	}
}

// drawPiece draws a single piece from the board at its grid position (i, j).
//
//go:noinline
func drawPiece(screen *ebiten.Image, b *board.Board, i, j, mouseX, mouseY int) {
	x := config.GridOriginX + i*(config.SquareSize+config.Gap)
	y := config.GridOriginY + j*(config.SquareSize+config.Gap)
	selectedX, selectedY := b.Selected()
	isSelected := (i == selectedX && j == selectedY)
	isHovered := mouseX >= x && mouseX < x+config.SquareSize && mouseY >= y && mouseY < y+config.SquareSize

	drawX, drawY := float64(x), float64(y)

	// Apply hover effect if not selected
	if isHovered && !isSelected {
		drawX += 1
		drawY += 1
	}

	// Draw shadow unless it's selected or hovered
	if !isSelected && !isHovered {
		shadowOffset := 2
		vector.DrawFilledRect(screen, float32(x+shadowOffset), float32(y+shadowOffset), float32(config.SquareSize), float32(config.SquareSize), config.ShadowColor, false)
	}

	color := b.Grid()[j][i]
	drawPieceAt(screen, color, drawX, drawY, isSelected)
}

// drawPieceAt dispatches drawing to specialized functions based on selection state.
//
//go:noinline
func drawPieceAt(screen *ebiten.Image, pieceColor color.Color, x, y float64, isSelected bool) {
	if pieceColor == nil {
		return // Don't draw empty cells
	}
	if isSelected {
		drawSelectedPiece(screen, pieceColor, x, y)
	} else {
		drawRegularPiece(screen, pieceColor, x, y)
	}
}

// drawRegularPiece draws a standard square piece.
//
//go:noinline
func drawRegularPiece(screen *ebiten.Image, pieceColor color.Color, x, y float64) {
	accentColor, ok := config.AccentColors[pieceColor]
	if !ok {
		accentColor = config.White // Default to white
	}

	square := ebiten.NewImage(config.SquareSize, config.SquareSize)
	square.Fill(pieceColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(square, op)

	accentSize := config.SquareSize / 4
	accentSquare := ebiten.NewImage(accentSize, accentSize)
	accentSquare.Fill(accentColor)
	accentOp := &ebiten.DrawImageOptions{}
	accentOp.GeoM.Translate(x+float64(config.SquareSize)/8, y+float64(config.SquareSize)/8)
	screen.DrawImage(accentSquare, accentOp)
}

// drawSelectedPiece draws a circular selected piece.
//
//go:noinline
func drawSelectedPiece(screen *ebiten.Image, pieceColor color.Color, x, y float64) {
	accentColor, ok := config.AccentColors[pieceColor]
	if !ok {
		accentColor = config.White // Default to white
	}

	shadowOffset := 2
	cx := float32(x + float64(config.SquareSize)/2)
	cy := float32(y + float64(config.SquareSize)/2)
	r := float32(config.SquareSize / 2)
	vector.DrawFilledCircle(screen, cx+float32(shadowOffset), cy+float32(shadowOffset), r, config.ShadowColor, true)
	vector.DrawFilledCircle(screen, cx, cy, r, pieceColor, true)
	accentR := float32(config.SquareSize / 8)
	vector.DrawFilledCircle(screen, cx, cy, accentR, accentColor, true)
}

//go:noinline
func drawUI(screen *ebiten.Image, score, maxScore, moveCount int, scoreHistory []int) {
	// The UI area is the space above the grid. We'll have a status bar and a graph area.
	const statusBarHeight = 30 // Height for the text labels at the top.
	uiSideMargin := 20

	// --- Draw text labels at the very top of the screen ---

	// Calculate the Y position to center the text line within the status bar area.
	// We get the text height once and use it for all labels.
	textBounds, _ := font.BoundString(config.STextFace, "M") // Get height of a sample character
	textH := (textBounds.Max.Y - textBounds.Min.Y).Ceil()
	textY := (statusBarHeight-textH)/2 + textH

	// Max Score (Top-Left)
	maxScoreStr := fmt.Sprintf("Max: %d", maxScore)
	text.Draw(screen, maxScoreStr, config.STextFace, uiSideMargin, textY, config.Black)

	// Current Score (Center-Left)
	scoreStr := fmt.Sprintf("Score: %d", score)
	scoreBounds, _ := font.BoundString(config.STextFace, scoreStr)
	scoreW := (scoreBounds.Max.X - scoreBounds.Min.X).Ceil()
	scoreX := (config.ScreenWidth - scoreW) / 2
	text.Draw(screen, scoreStr, config.STextFace, scoreX, textY, config.Black)

	// Move Counter (Bottom-Right)
	moveCountStr := fmt.Sprintf("Moves: %d", moveCount)
	moveBounds, _ := font.BoundString(config.STextFace, moveCountStr)
	moveW := (moveBounds.Max.X - moveBounds.Min.X).Ceil()
	moveX := config.ScreenWidth - moveW - uiSideMargin
	text.Draw(screen, moveCountStr, config.STextFace, moveX, textY, config.Black)

	// --- Draw the score graph below the status bar ---
	graphX := 0
	graphY := statusBarHeight
	graphWidth := config.ScreenWidth
	const graphBottomMargin = 20 // Space between graph and board
	graphHeight := config.GridOriginY - statusBarHeight - graphBottomMargin
	drawScoreGraph(screen, scoreHistory, maxScore, graphX, graphY, graphWidth, graphHeight)
}

//go:noinline
func drawScoreGraph(screen *ebiten.Image, history []int, maxScore, x, y, width, height int) {
	// Draw graph background/border
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), config.LightGrey, false)

	if len(history) < 2 {
		return // Not enough data to draw a line
	}

	// Determine the Y-axis scale. Use max possible score to prevent the scale from jumping around.
	// If the current score exceeds the theoretical max (e.g., due to future scoring changes), adjust.
	yMax := float32(maxScore)
	currentMax := 0
	for _, s := range history {
		if s > currentMax {
			currentMax = s
		}
	}
	if float32(currentMax) > yMax {
		yMax = float32(currentMax)
	}
	if yMax == 0 {
		yMax = 1 // Avoid division by zero
	}

	// Determine the X-axis scale. The graph will show up to a certain number of recent moves,
	// but we will scale the drawing to fit the full width.
	maxVisibleMoves := width // Show up to `width` moves on the graph
	historyLen := len(history)
	startIdx := 0
	if historyLen > maxVisibleMoves {
		startIdx = historyLen - maxVisibleMoves
	}
	visibleHistory := history[startIdx:]
	numVisiblePoints := len(visibleHistory)

	// Keep track of the last point where the score was not zero, for connecting lines.
	lastY := float32(y + height) // Start at the baseline
	lastX := float32(x)

	// Draw the line segments
	for i := 0; i < numVisiblePoints; i++ {
		currentScore := float32(visibleHistory[i])
		previousScore := float32(0)
		if i > 0 {
			previousScore = float32(visibleHistory[i-1])
		}

		// The current X position is based on the index in the visible history. 1 move = 1 pixel.
		currentX := float32(x + i)
		currentY := float32(y+height) - (currentScore/yMax)*float32(height)

		if currentScore > previousScore {
			// Score increased: draw a green line from the last point to the current one.
			vector.StrokeLine(screen, lastX, lastY, currentX, currentY, 2, config.Green, true)
		} else if currentScore < previousScore {
			// Score decreased: draw a red line from the last point to the current one.
			vector.StrokeLine(screen, lastX, lastY, currentX, currentY, 2, config.Red, true)
		} else if i > 0 { // Only draw dots for moves after the first one.
			// Score is unchanged: draw a subtle grey line segment.
			// This creates a continuous horizontal line.
			vector.StrokeLine(screen, lastX, lastY, currentX, lastY, 1, config.Grey, true)
		}

		// Update the last position for the next iteration's line segment.
		lastX = currentX
		lastY = currentY
	}
}

//go:noinline
func drawStoneDistribution(screen *ebiten.Image, counts map[color.Color]int) {
	// A struct to hold color and count for sorting
	type colorCount struct {
		Color color.Color
		Count int
	}

	// Helper function to convert a color to a comparable string for stable sorting
	colorToString := func(c color.Color) string {
		r, g, b, a := c.RGBA()
		// Use a format that ensures lexicographical sorting matches color value sorting
		return fmt.Sprintf("%05d-%05d-%05d-%05d", r, g, b, a)
	}

	// Convert map to a slice for stable, sorted display
	var sortedCounts []colorCount
	for c, n := range counts {
		sortedCounts = append(sortedCounts, colorCount{Color: c, Count: n})
	}

	// Sort by count (descending), then by color (ascending) for a stable order.
	// This prevents flickering when counts are equal.
	sort.Slice(sortedCounts, func(i, j int) bool {
		if sortedCounts[i].Count != sortedCounts[j].Count {
			return sortedCounts[i].Count > sortedCounts[j].Count
		}
		return colorToString(sortedCounts[i].Color) < colorToString(sortedCounts[j].Color)
	})

	// --- Drawing constants ---
	// Define the area below the board for the miniatures.
	areaTopY := config.GridOriginY + config.GridHeight
	areaBottomY := config.ScreenHeight
	areaHeight := areaBottomY - areaTopY

	miniatureSize := 24
	itemHeight := 30
	textMarginLeft := 10
	numColumns := 4
	itemWidth := 80 // Width for one "icon + text" block

	// Center the block of miniatures both horizontally and vertically in the space below the board.
	blockStartX := (config.ScreenWidth - (numColumns * itemWidth)) / 2
	blockStartY := areaTopY + (areaHeight-((len(sortedCounts)+numColumns-1)/numColumns*itemHeight))/2

	for i, item := range sortedCounts {
		col := i % numColumns
		row := i / numColumns

		x := blockStartX + col*itemWidth
		y := blockStartY + row*itemHeight

		// Draw the color miniature
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(miniatureSize), float32(miniatureSize), item.Color, true)

		// Draw the accent color on top of the miniature
		accentColor, ok := config.AccentColors[item.Color]
		if !ok {
			accentColor = config.White // Default accent
		}
		accentSize := float32(miniatureSize / 4)
		accentOffset := float32(miniatureSize / 8)
		accentX := float32(x) + accentOffset
		accentY := float32(y) + accentOffset

		vector.DrawFilledRect(screen, accentX, accentY, accentSize, accentSize, accentColor, true)

		// Draw the count text
		countStr := fmt.Sprintf("%d", item.Count)
		textBounds, _ := font.BoundString(config.MTextFace, countStr)
		textH := (textBounds.Max.Y - textBounds.Min.Y).Ceil()
		textX := x + miniatureSize + textMarginLeft
		// Vertically center the text next to the miniature
		textY := y + (miniatureSize-textH)/2 + textH - 2

		text.Draw(screen, countStr, config.MTextFace, textX, textY, config.Black)
	}
}
