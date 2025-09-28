package view

import (
	"fmt"
	"image/color"
	"zenmojo/board"
	"zenmojo/config"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw renders the entire game screen.
//
//go:noinline
func Draw(screen *ebiten.Image, b *board.Board, score, maxScore, moveCount, mouseX, mouseY int) {
	drawBackground(screen)
	drawBoard(screen, b, mouseX, mouseY)
	drawUI(screen, score, maxScore, moveCount)
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
func drawUI(screen *ebiten.Image, score, maxScore, moveCount int) {
	uiTopMargin := 40
	uiSideMargin := 20

	// Draw move counter on the top left
	moveCountStr := fmt.Sprintf("Moves: %d", moveCount)
	text.Draw(screen, moveCountStr, config.MTextFace, uiSideMargin, uiTopMargin, config.Black)

	// --- Draw Score Progress Bar on the top right ---

	// Define bar dimensions and position
	barWidth := 200
	barHeight := 20
	barX := config.ScreenWidth - barWidth - uiSideMargin
	barY := uiTopMargin - barHeight/2 // Align vertically with "Moves" text

	// Calculate progress
	progress := 0.0
	if maxScore > 0 {
		progress = float64(score) / float64(maxScore)
	}
	if progress > 1.0 { // Cap progress at 100%
		progress = 1.0
	}

	// Draw bar background
	vector.DrawFilledRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), config.LightGrey, false)

	// Draw filled portion of the bar
	fillWidth := float32(float64(barWidth) * progress)
	vector.DrawFilledRect(screen, float32(barX), float32(barY), fillWidth, float32(barHeight), config.Green, false)

	// Draw progress text (e.g., "123 / 456") over the bar
	progressStr := fmt.Sprintf("%d / %d", score, maxScore)
	textBounds := text.BoundString(config.MTextFace, progressStr)
	textX := barX + (barWidth-textBounds.Dx())/2
	textY := barY + (barHeight+textBounds.Dy())/2 - 4 // Adjust for vertical centering
	text.Draw(screen, progressStr, config.MTextFace, textX, textY, config.White)
}
