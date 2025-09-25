package config

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 640
	GridSize     = 10
	SquareSize   = 48
	Gap          = 8
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

	Icons []image.Image
)

func init() {
	// Calculate grid dimensions
	GridWidth = GridSize*SquareSize + (GridSize-1)*Gap
	GridHeight = GridSize*SquareSize + (GridSize-1)*Gap
	GridOriginX = (ScreenWidth - GridWidth) / 2
	GridOriginY = (ScreenHeight - GridHeight) / 2

	// Create graphical assets
	patternSize := 20
	HatchingPattern = ebiten.NewImage(patternSize, patternSize)
	for i := -patternSize; i < patternSize; i += 4 {
		vector.StrokeLine(HatchingPattern, float32(i), 0, float32(i+patternSize), float32(patternSize), 1, HatchingColor, false)
	}

	// Load font
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

	// Create window icons
	Icons = []image.Image{createRainbowIcon(16), createRainbowIcon(32), createRainbowIcon(48)}
}

// createRainbowIcon is a helper function to generate the window icon.
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
