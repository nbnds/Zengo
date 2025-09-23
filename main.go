package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 640
	GridSize     = 10
	squareSize   = 48
	gap          = 8
)

var (
	// A curated color palette
	Red        = color.RGBA{R: 237, G: 63, B: 39, A: 255}
	OldRed     = color.RGBA{R: 171, G: 68, B: 89, A: 255}
	Yellow     = color.RGBA{R: 255, G: 204, B: 0, A: 255}
	Green      = color.RGBA{R: 52, G: 199, B: 89, A: 255}
	Blue       = color.RGBA{R: 0, G: 122, B: 255, A: 255}
	Black      = color.RGBA{R: 27, G: 24, B: 51, A: 255}
	LightBrown = color.RGBA{R: 210, G: 180, B: 140, A: 255}
	White      = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	CrimsonRed = color.RGBA{R: 220, G: 20, B: 60, A: 255}
	Plume      = color.RGBA{R: 107, G: 63, B: 105, A: 255}
	Grey       = color.RGBA{R: 55, G: 53, B: 62, A: 255}
	LightGrey  = color.RGBA{R: 211, G: 218, B: 217, A: 255}
	Peach      = color.RGBA{R: 242, G: 159, B: 88, A: 255}
	Orange     = color.RGBA{R: 255, G: 149, B: 0, A: 255}
	Purple     = color.RGBA{R: 175, G: 82, B: 222, A: 255}
	Cyan       = color.RGBA{R: 50, G: 173, B: 230, A: 255}
	Magenta    = color.RGBA{R: 255, G: 45, B: 85, A: 255}
	Lime       = color.RGBA{R: 204, G: 255, B: 0, A: 255}
	Teal       = color.RGBA{R: 90, G: 200, B: 250, A: 255}
	Brown      = color.RGBA{R: 162, G: 132, B: 94, A: 255}
	Pink       = color.RGBA{R: 255, G: 105, B: 180, A: 255}
	Gold       = color.RGBA{R: 255, G: 215, B: 0, A: 255}
	Silver     = color.RGBA{R: 192, G: 192, B: 192, A: 255}
	DarkGreen  = color.RGBA{R: 10, G: 64, B: 12, A: 255}

	Palette = []color.Color{
		Red, Yellow, Green, Blue, Black, LightBrown, Plume, Grey, Peach, OldRed,
		Orange, Purple, Cyan, Magenta, Lime, Teal, Brown, Pink, Gold, Silver, DarkGreen,
	}

	AccentColors = map[color.Color]color.Color{
		Red:        Orange,
		Yellow:     Red,
		Green:      Black,
		Blue:       White,
		Black:      CrimsonRed,
		LightBrown: Black,
		Plume:      White,
		Grey:       LightGrey,
		Peach:      OldRed,
		OldRed:     Peach,
		Orange:     White,
		Purple:     White,
		Cyan:       Black,
		Magenta:    White,
		Lime:       Black,
		Teal:       Black,
		Brown:      White,
		Pink:       Black,
		Gold:       Black,
		Silver:     Red,
		DarkGreen:  White,
	}

	GridWidth   int
	GridHeight  int
	GridOriginX int
	GridOriginY int

	BackgroundColor = color.RGBA{R: 245, G: 239, B: 230, A: 255}
	HatchingColor   = color.RGBA{R: 203, G: 220, B: 235, A: 255} // Light purple
	ShadowColor     = color.RGBA{R: 0, G: 0, B: 0, A: 128}

	HatchingPattern *ebiten.Image

	MTextFace font.Face
)

func init() {
	GridWidth = GridSize*squareSize + (GridSize-1)*gap
	GridHeight = GridSize*squareSize + (GridSize-1)*gap
	GridOriginX = (ScreenWidth - GridWidth) / 2
	GridOriginY = (ScreenHeight - GridHeight) / 2

	patternSize := 20
	HatchingPattern = ebiten.NewImage(patternSize, patternSize)
	for i := -patternSize; i < patternSize; i += 4 {
		vector.StrokeLine(HatchingPattern, float32(i), 0, float32(i+patternSize), float32(patternSize), 1, HatchingColor, false)
	}

	ttf, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	baseFace, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	MTextFace = baseFace
}

type Game struct {
	mouseX    int
	mouseY    int
	selectedX int
	selectedY int
	grid      [][]color.Color
	moveCount int
	score     int // Added score field
}

// Coordinate represents a position on the grid
type Coordinate struct {
	R, C int
}

// Group represents a contiguous group of same-colored items
type Group struct {
	Color       color.Color
	Coordinates []Coordinate
	MinR, MaxR  int // Bounding box for shape detection
	MinC, MaxC  int // Bounding box for shape detection
}

func (g *Game) findGroups() []Group {
	visited := make([][]bool, GridSize)
	for i := range visited {
		visited[i] = make([]bool, GridSize)
	}

	var groups []Group

	for r := 0; r < GridSize; r++ {
		for c := 0; c < GridSize; c++ {
			if !visited[r][c] {
				currentGroup := Group{
					Color:       g.grid[r][c],
					Coordinates: []Coordinate{},
					MinR:        r,
					MaxR:        r,
					MinC:        c,
					MaxC:        c,
				}

				g.dfs(r, c, g.grid[r][c], &currentGroup, visited)
				groups = append(groups, currentGroup)
			}
		}
	}
	return groups
}

func (g *Game) dfs(r, c int, targetColor color.Color, currentGroup *Group, visited [][]bool) {
	// Check bounds and if already visited or color mismatch
	if r < 0 || r >= GridSize || c < 0 || c >= GridSize || visited[r][c] || !colorsEqual(g.grid[r][c], targetColor) {
		return
	}

	visited[r][c] = true
	currentGroup.Coordinates = append(currentGroup.Coordinates, Coordinate{R: r, C: c})

	// Update bounding box
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

	// Explore neighbors
	g.dfs(r+1, c, targetColor, currentGroup, visited)
	g.dfs(r-1, c, targetColor, currentGroup, visited)
	g.dfs(r, c+1, targetColor, currentGroup, visited)
	g.dfs(r, c-1, targetColor, currentGroup, visited)
}

// colorsEqual is a helper function to compare two color.Color values
func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func (g *Game) calculateScore(group Group) int {
	// A group must have at least 2 items to be considered for scoring
	if len(group.Coordinates) < 2 {
		return 0
	}

	width := group.MaxC - group.MinC + 1
	height := group.MaxR - group.MinR + 1

	// Line shape (for example 7 items) we count one point for each item in the group
	if width == 1 || height == 1 {
		return len(group.Coordinates)
	}

	// 2 by X shape (for example 6 items arrranged in 2x3, we count 6 (number of items in the group) multiplied by the one side (2) and by the other side (3) = 6 x 6 = 36
	if (width == 2 && height >= 2) || (height == 2 && width >= 2) {
		// Ensure it's a perfect rectangle of 2xX or Xx2
		if len(group.Coordinates) == width*height {
			return len(group.Coordinates) * width * height
		}
	}

	// 3x3 shape (9 items) we count 9 points for group by 3 by 3 = 9 x 9 = 81 points.
	if width == 3 && height == 3 {
		// Ensure it's a perfect 3x3 square
		if len(group.Coordinates) == 9 {
			return len(group.Coordinates) * width * height
		}
	}

	return 0 // No score for other shapes or incomplete rectangles
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		clicked := false
		for i := 0; i < GridSize; i++ {
			for j := 0; j < GridSize; j++ {
				x := GridOriginX + i*(squareSize+gap)
				y := GridOriginY + j*(squareSize+gap)
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

						// Calculate score after swap
						currentScore := 0
						groups := g.findGroups()
						for _, group := range groups {
							currentScore += g.calculateScore(group)
						}
						g.score = currentScore

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
	screen.Fill(BackgroundColor)

	// Draw hatching pattern
	patternSize := HatchingPattern.Bounds().Dx()
	for i := 0; i < ScreenWidth; i += patternSize {
		for j := 0; j < ScreenHeight; j += patternSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i), float64(j))
			screen.DrawImage(HatchingPattern, op)
		}
	}

	shadowOffset := 2

	for i := range GridSize {
		for j := range GridSize {
			x := GridOriginX + i*(squareSize+gap)
			y := GridOriginY + j*(squareSize+gap)

			hovered := g.mouseX >= x && g.mouseX < x+squareSize && g.mouseY >= y && g.mouseY < y+squareSize

			drawX, drawY := x, y

			// Draw shadow
			if !hovered && !(i == g.selectedX && j == g.selectedY) {
				vector.DrawFilledRect(screen, float32(x+shadowOffset), float32(y+shadowOffset), float32(squareSize), float32(squareSize), ShadowColor, false)
			} else if hovered {
				drawX += 1
				drawY += 1
			}

			mainColor := g.grid[j][i]
			accentColor, ok := AccentColors[mainColor]
			if !ok {
				accentColor = White // Default to white
			}

			if i == g.selectedX && j == g.selectedY {
				// Draw circular shadow for selected circle
				cx := float32(drawX + squareSize/2)
				cy := float32(drawY + squareSize/2)
				r := float32(squareSize / 2)
				vector.DrawFilledCircle(screen, cx+float32(shadowOffset), cy+float32(shadowOffset), r, ShadowColor, true)

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
	text.Draw(screen, moveCountStr, MTextFace, 10, ScreenHeight-10, color.Black)

	// Draw score
	scoreStr := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreStr, MTextFace, 10, ScreenHeight-40, color.Black)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func CreateRainbowIcon(size int) image.Image {
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

func generateColorCounts(palette []color.Color, GridSize int) []int {
	counts := make([]int, len(palette))
	sum := 0

	// First, assign a random number between 2 and 10 to each color
	for i := range palette {
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
	icons := []image.Image{CreateRainbowIcon(16), CreateRainbowIcon(32), CreateRainbowIcon(48)}
	ebiten.SetWindowIcon(icons)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Zesty Zen")

	game := &Game{
		selectedX: -1,
		selectedY: -1,
		grid:      make([][]color.Color, GridSize),
		moveCount: 0,
	}

	counts := generateColorCounts(Palette, GridSize)

	shuffledColors := make([]color.Color, 0, GridSize*GridSize)
	for i, c := range Palette {
		for j := 0; j < counts[i]; j++ {
			shuffledColors = append(shuffledColors, c)
		}
	}

	rand.Shuffle(len(shuffledColors), func(i, j int) {
		shuffledColors[i], shuffledColors[j] = shuffledColors[j], shuffledColors[i]
	})

	// Populate the grid with shuffled colors
	for i := range game.grid {
		game.grid[i] = make([]color.Color, GridSize)
		for j := range game.grid[i] {
			game.grid[i][j] = shuffledColors[i*GridSize+j]
		}
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
