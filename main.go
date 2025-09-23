package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 800
	screenHeight = 640
	gridSize     = 10
	squareSize   = 48
	gap          = 8
)

var (
	// A curated color palette
	red        = color.RGBA{R: 237, G: 63, B: 39, A: 255}
	oldRed     = color.RGBA{R: 171, G: 68, B: 89, A: 255}
	yellow     = color.RGBA{R: 255, G: 204, B: 0, A: 255}
	green      = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	blue       = color.RGBA{R: 0, G: 122, B: 255, A: 255}
	black      = color.RGBA{R: 27, G: 24, B: 51, A: 255}
	lightBrown = color.RGBA{R: 210, G: 180, B: 140, A: 255}
	white      = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	crimsonRed = color.RGBA{R: 220, G: 20, B: 60, A: 255}
	plume      = color.RGBA{R: 107, G: 63, B: 105, A: 255}
	grey       = color.RGBA{R: 55, G: 53, B: 62, A: 255}
	lightGrey  = color.RGBA{R: 211, G: 218, B: 217, A: 255}
	peach      = color.RGBA{R: 242, G: 159, B: 88, A: 255}
	orange     = color.RGBA{R: 255, G: 149, B: 0, A: 255}
	purple     = color.RGBA{R: 175, G: 82, B: 222, A: 255}
	cyan       = color.RGBA{R: 50, G: 173, B: 230, A: 255}
	magenta    = color.RGBA{R: 255, G: 45, B: 85, A: 255}
	lime       = color.RGBA{R: 204, G: 255, B: 0, A: 255}
	teal       = color.RGBA{R: 90, G: 200, B: 250, A: 255}
	brown      = color.RGBA{R: 162, G: 132, B: 94, A: 255}
	pink       = color.RGBA{R: 255, G: 105, B: 180, A: 255}
	gold       = color.RGBA{R: 255, G: 215, B: 0, A: 255}
	silver     = color.RGBA{R: 192, G: 192, B: 192, A: 255}
	darkGreen  = color.RGBA{R: 10, G: 64, B: 12, A: 255}

	palette = []color.Color{
		red, yellow, green, blue, black, lightBrown, plume, grey, peach, oldRed,
		orange, purple, cyan, magenta, lime, teal, brown, pink, gold, silver, darkGreen,
	}

	accentColors = map[color.Color]color.Color{
		red:        orange,
		yellow:     red,
		green:      black,
		blue:       white,
		black:      crimsonRed,
		lightBrown: black,
		plume:      white,
		grey:       lightGrey,
		peach:      oldRed,
		oldRed:     peach,
		orange:     white,
		purple:     white,
		cyan:       black,
		magenta:    white,
		lime:       black,
		teal:       black,
		brown:      white,
		pink:       black,
		gold:       black,
		silver:     red,
		darkGreen:  white,
	}

	gridWidth   int
	gridHeight  int
	gridOriginX int
	gridOriginY int

	backgroundColor = color.RGBA{R: 245, G: 239, B: 230, A: 255}
	hatchingColor   = color.RGBA{R: 203, G: 220, B: 235, A: 255} // Light purple
	shadowColor     = color.RGBA{R: 0, G: 0, B: 0, A: 128}

	hatchingPattern *ebiten.Image

	mTextFace font.Face
)

func init() {
	gridWidth = gridSize*squareSize + (gridSize-1)*gap
	gridHeight = gridSize*squareSize + (gridSize-1)*gap
	gridOriginX = (screenWidth - gridWidth) / 2
	gridOriginY = (screenHeight - gridHeight) / 2

	patternSize := 20
	hatchingPattern = ebiten.NewImage(patternSize, patternSize)
	for i := -patternSize; i < patternSize; i += 4 {
		vector.StrokeLine(hatchingPattern, float32(i), 0, float32(i+patternSize), float32(patternSize), 1, hatchingColor, false)
	}

	ttf, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mTextFace, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	mouseX    int
	mouseY    int
	selectedX int
	selectedY int
	grid      [][]color.Color
	moveCount int
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		clicked := false
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				x := gridOriginX + i*(squareSize+gap)
				y := gridOriginY + j*(squareSize+gap)
				if g.mouseX >= x && g.mouseX < x+squareSize && g.mouseY >= y && g.mouseY < y+squareSize {
					if g.selectedX == -1 {
						// Select a square
						g.selectedX = i
						g.selectedY = j
					} else if g.selectedX == i && g.selectedY == j {
						// Deselect if clicking the same square
						g.selectedX = -1
						g.selectedY = -1
					} else {
						// Swap squares
						g.grid[g.selectedY][g.selectedX], g.grid[j][i] = g.grid[j][i], g.grid[g.selectedY][g.selectedX]
						g.moveCount++
						g.selectedX = -1
						g.selectedY = -1
					}
					clicked = true
					break
				}
			}
			if clicked {
				break
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	// Draw hatching pattern
	patternSize := hatchingPattern.Bounds().Dx()
	for i := 0; i < screenWidth; i += patternSize {
		for j := 0; j < screenHeight; j += patternSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i), float64(j))
			screen.DrawImage(hatchingPattern, op)
		}
	}

	shadowOffset := 2

	for i := range gridSize {
		for j := range gridSize {
			x := gridOriginX + i*(squareSize+gap)
			y := gridOriginY + j*(squareSize+gap)

			hovered := g.mouseX >= x && g.mouseX < x+squareSize && g.mouseY >= y && g.mouseY < y+squareSize

			drawX, drawY := x, y

			// Draw shadow
			if !hovered && !(i == g.selectedX && j == g.selectedY) {
				vector.DrawFilledRect(screen, float32(x+shadowOffset), float32(y+shadowOffset), float32(squareSize), float32(squareSize), shadowColor, false)
			} else if hovered {
				drawX += 1
				drawY += 1
			}

			mainColor := g.grid[j][i]
			accentColor, ok := accentColors[mainColor]
			if !ok {
				accentColor = white // Default to white
			}

			if i == g.selectedX && j == g.selectedY {
				// Draw circular shadow for selected circle
				cx := float32(drawX + squareSize/2)
				cy := float32(drawY + squareSize/2)
				r := float32(squareSize / 2)
				vector.DrawFilledCircle(screen, cx+float32(shadowOffset), cy+float32(shadowOffset), r, shadowColor, true)

				// Draw selected circle
				vector.DrawFilledCircle(screen, cx, cy, r, mainColor, true)

				// Draw accent circle in the middle
				accentR := float32(squareSize / 8)
				vector.DrawFilledCircle(screen, cx, cy, accentR, accentColor, true)

			} else {
				// Draw normal square
				square := ebiten.NewImage(squareSize, squareSize)
				square.Fill(mainColor)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(drawX), float64(drawY))
				screen.DrawImage(square, op)

				// Draw accent square
				accentSize := squareSize / 4
				accentSquare := ebiten.NewImage(accentSize, accentSize)
				accentSquare.Fill(accentColor)
				accentOp := &ebiten.DrawImageOptions{}
				accentOp.GeoM.Translate(float64(drawX+squareSize/8), float64(drawY+squareSize/8))
				screen.DrawImage(accentSquare, accentOp)
			}
		}
	}

	// Draw move counter
	moveCountStr := fmt.Sprintf("Moves: %d", g.moveCount)
	text.Draw(screen, moveCountStr, mTextFace, 10, screenHeight-10, color.Black)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func createRainbowIcon(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / size),
				G: uint8(y * 255 / size),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func generateColorCounts(palette []color.Color, gridSize int) []int {
	rand.Seed(time.Now().UnixNano())
	counts := make([]int, len(palette))
	sum := 0

	// First, assign a random number between 2 and 10 to each color
	for i := 0; i < len(palette); i++ {
		counts[i] = rand.Intn(9) + 2 // 2 to 10
	}

	// Now, adjust the counts to make their sum 100
	for {
		sum = 0
		for _, c := range counts {
			sum += c
		}
		if sum == 100 {
			break
		}
		// Pick a random color and increment or decrement its count
		idx := rand.Intn(len(palette))
		if sum > 100 {
			if counts[idx] > 2 {
				counts[idx]--
			}
		} else {
			if counts[idx] < 10 {
				counts[idx]++
			}
		}
	}
	return counts
}

func main() {
	icons := []image.Image{createRainbowIcon(16), createRainbowIcon(32), createRainbowIcon(48)}
	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Zesty Zen")

	game := &Game{
		selectedX: -1,
		selectedY: -1,
		grid:      make([][]color.Color, gridSize),
		moveCount: 0,
	}

	counts := generateColorCounts(palette, gridSize)

	shuffledColors := make([]color.Color, 0, gridSize*gridSize)
	for i, c := range palette {
		for j := 0; j < counts[i]; j++ {
			shuffledColors = append(shuffledColors, c)
		}
	}

	rand.Shuffle(len(shuffledColors), func(i, j int) {
		shuffledColors[i], shuffledColors[j] = shuffledColors[j], shuffledColors[i]
	})

	// Populate the grid with shuffled colors
	for i := range game.grid {
		game.grid[i] = make([]color.Color, gridSize)
		for j := range game.grid[i] {
			game.grid[i][j] = shuffledColors[i*gridSize+j]
		}
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
