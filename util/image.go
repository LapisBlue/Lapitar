package util

import "image"

func IsSolidColor(img image.Image) bool {
	base := img.At(img.Bounds().Min.X, img.Bounds().Min.Y)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if img.At(x, y) != base {
				return false
			}
		}
	}

	return true
}
