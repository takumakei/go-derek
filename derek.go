// Package derek implements Derek Bradley's "Adaptive Thresholding using the Integral Image".
package derek

import (
	"image"
)

// ProcessGray returns the image with the Derek Bradley's "Adaptive Thresholding using the Integral Image" filter applied.
// threshold must be between [0, 100].
func ProcessGray(src *image.Gray, clusterSize int, threshold int) *image.Gray {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return src
	}

	si := src.Pix
	sw := src.Stride

	ii := make([]uint, w*h)
	for y := 0; y < h; y++ {
		var sum uint
		s := y * sw
		i := y * w
		for x := 0; x < w; x++ {
			sum += uint(si[x+y*sw])
			if y == 0 {
				ii[i] = sum
			} else {
				ii[i] = sum + ii[i-w]
			}
			s++
			i++
		}
	}

	thresholdu := uint(min(max(threshold, 0), 100))
	clusterSize = max(clusterSize, 0) + 1
	dst := image.NewGray(bounds)
	pix := dst.Pix
	piw := dst.Stride
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			x1 := max(x-clusterSize, 0)
			x2 := min(x+clusterSize, w-1)
			y1 := max(y-clusterSize, 0)
			y2 := min(y+clusterSize, h-1)
			count := (uint(x2) - uint(x1)) * (uint(y2) - uint(y1))
			y1w := y1 * w
			y2w := y2 * w
			sum := ii[x2+y2w] - ii[x1+y2w] - ii[x2+y1w] + ii[x1+y1w]
			pxv := uint(si[x+y*sw])
			if pxv*count*100 > sum*thresholdu {
				pix[x+y*piw] = 255
			}
		}
	}

	return dst
}

// Process returns the image with the Derek Bradley's "Adaptive Thresholding using the Integral Image" filter applied.
// threshold must be between [0, 100].
func Process(src image.Image, clusterSize int, threshold int) image.Image {
	return ProcessGray(NewGray(src), clusterSize, threshold)
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
