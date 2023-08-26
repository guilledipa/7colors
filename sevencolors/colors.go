package sevencolors

import (
	"image/color"
	"math/rand"
)

func tileBackgroundColor(value int) color.Color {
	switch value {
	case 0:
		return color.RGBA{0xff, 0x00, 0x00, 0xff} // Red
	case 1:
		return color.RGBA{0x00, 0xff, 0x00, 0xff} // Green
	case 2:
		return color.RGBA{0x00, 0x00, 0xff, 0xff} // Blue
	case 3:
		return color.RGBA{0xff, 0xff, 0x00, 0xff} // Yellow
	case 4:
		return color.RGBA{0xff, 0x00, 0xff, 0xff} // Magenta
	case 5:
		return color.RGBA{0x00, 0xff, 0xff, 0xff} // Cyan
	case 6:
		return color.RGBA{0xee, 0xee, 0xee, 0xff} // Grey
	case 7:
		return color.RGBA{0x00, 0x00, 0x00, 0xff} // Black
	}
	panic("not reach")
}

func generateRandomColor(rng *rand.Rand) color.Color {
	v := rng.Intn(7)
	return tileBackgroundColor(v)
}
