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
func Draw(screen *ebiten.Image, b *board.Board, score int, moveCount int, mouseX int, mouseY int) {
	drawBackground(screen)
	drawBoard(screen, b, mouseX, mouseY)
	drawUI(screen, score, moveCount)
}

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

func drawBoard(screen *ebiten.Image, b *board.Board, mouseX int, mouseY int) {
	shadowOffset := 2
	grid := b.Grid()
	selectedX, selectedY := b.Selected()

	for i := range config.GridSize {
		for j := range config.GridSize {
			x := config.GridOriginX + i*(config.SquareSize+config.Gap)
			y := config.GridOriginY + j*(config.SquareSize+config.Gap)

			hovered := mouseX >= x && mouseX < x+config.SquareSize && mouseY >= y && mouseY < y+config.SquareSize

			drawX, drawY := x, y

			// Draw shadow
			if !hovered && !(i == selectedX && j == selectedY) {
				vector.DrawFilledRect(screen, float32(x+shadowOffset), float32(y+shadowOffset), float32(config.SquareSize), float32(config.SquareSize), config.ShadowColor, false)
			} else if hovered {
				drawX += 1
				drawY += 1
			}

			mainColor := grid[j][i]
			if mainColor == nil {
				continue // Don't draw empty cells
			}
			accentColor, ok := config.AccentColors[mainColor]
			if !ok {
				accentColor = config.White // Default to white
			}

			if i == selectedX && j == selectedY {
				// Draw circular shadow for selected circle
				cx := float32(drawX + config.SquareSize/2)
				cy := float32(drawY + config.SquareSize/2)
				r := float32(config.SquareSize / 2)
				vector.DrawFilledCircle(screen, cx+float32(shadowOffset), cy+float32(shadowOffset), r, config.ShadowColor, true)

				// Draw selected circle
				vector.DrawFilledCircle(screen, cx, cy, r, mainColor, true)

				// Draw accent circle in the middle
				accentR := float32(config.SquareSize / 8)
				vector.DrawFilledCircle(screen, cx, cy, accentR, accentColor, true)

			} else {
				// Draw normal square
				square := ebiten.NewImage(config.SquareSize, config.SquareSize)
				square.Fill(mainColor)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(drawX), float64(drawY))
				screen.DrawImage(square, op)

				// Draw accent square
				accentSize := config.SquareSize / 4
				accentSquare := ebiten.NewImage(accentSize, accentSize)
				accentSquare.Fill(accentColor)
				accentOp := &ebiten.DrawImageOptions{}
				accentOp.GeoM.Translate(float64(drawX+config.SquareSize/8), float64(drawY+config.SquareSize/8))
				screen.DrawImage(accentSquare, accentOp)
			}
		}
	}
}

func drawUI(screen *ebiten.Image, score int, moveCount int) {
	// Draw move counter
	moveCountStr := fmt.Sprintf("Moves: %d", moveCount)
	text.Draw(screen, moveCountStr, config.MTextFace, 10, config.ScreenHeight-10, color.Black)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", score)
	text.Draw(screen, scoreStr, config.MTextFace, 10, config.ScreenHeight-40, color.Black)
}
