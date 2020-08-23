package derek

import (
	"image"
	"image/color"
)

// NewGray returns a grayscale image of src.
func NewGray(src image.Image) *image.Gray {
	if gray, ok := src.(*image.Gray); ok {
		return gray
	}

	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	gray := image.NewGray(b)
	graypx := gray.Pix
	stride := gray.Stride
	sy := b.Min.Y
	for y := 0; y < h; y++ {
		sx := b.Min.X
		di := y * stride
		for x := 0; x < w; x++ {
			graypx[di] = color.GrayModel.Convert(src.At(sx, sy)).(color.Gray).Y
			sx++
			di++
		}
		sy++
	}
	return gray
}
