// Package derek implements Derek Bradley's "Adaptive Thresholding using the Integral Image".
package derek

import (
	"image"
	"image/color"
)

// Process returns the image with the Derek Bradley's "Adaptive Thresholding using the Integral Image" filter applied.
// threshold must be between [0, 100].
func Process(src image.Image, clusterSize int, threshold int) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return src
	}

	gr := color.Gray16Model
	si := make([]uint, w*h)
	ii := make([]uint, w*h)
	sx := bounds.Min.X
	for x := 0; x < w; x++ {
		var sum uint
		sy := bounds.Min.Y
		for y := 0; y < h; y++ {
			i := y*w + x
			pxv := uint(gr.Convert(src.At(sx, sy)).(color.Gray16).Y)
			si[i] = pxv
			sum += pxv
			if x == 0 {
				ii[i] = sum
			} else {
				ii[i] = sum + ii[i-1]
			}
			sy++
		}
		sx++
	}

	clusterSize++
	dst := image.NewGray(image.Rect(0, 0, w, h))
	pix := dst.Pix
	piw := dst.Stride
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			x1 := max(x-clusterSize, 0)
			x2 := min(x+clusterSize, w-1)
			y1 := max(y-clusterSize, 0)
			y2 := min(y+clusterSize, h-1)
			count := (uint(x2) - uint(x1)) * (uint(y2) - uint(y1))
			y1w := y1 * w
			y2w := y2 * w
			sum := ii[x2+y2w] - ii[x1+y2w] - ii[x2+y1w] + ii[x1+y1w]
			rgb := si[x+y*w]
			if rgb*count > sum*uint(threshold)/100 {
				pix[x+y*piw] = 255
			}
		}
	}

	return dst
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
