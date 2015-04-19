package util

import (
	"errors"
	"github.com/disintegration/imaging"
	"image"
)

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

func ParseScale(name string) (filter *imaging.ResampleFilter, err error) {
	switch name {
	case "NearestNeighbor":
		filter = &imaging.NearestNeighbor
	case "Box":
		filter = &imaging.Box
	case "Linear":
		filter = &imaging.Linear
	case "Hermite":
		filter = &imaging.Hermite
	case "MitchellNetravali":
		filter = &imaging.MitchellNetravali
	case "CatmullRom":
		filter = &imaging.CatmullRom
	case "BSpline":
		filter = &imaging.BSpline
	case "Gaussian":
		filter = &imaging.Gaussian
	case "Bartlett":
		filter = &imaging.Bartlett
	case "Lanczos":
		filter = &imaging.Lanczos
	case "Hann":
		filter = &imaging.Hann
	case "Hamming":
		filter = &imaging.Hamming
	case "Blackman":
		filter = &imaging.Blackman
	case "Welch":
		filter = &imaging.Welch
	case "Cosine":
		filter = &imaging.Cosine
	default:
		err = errors.New("Unknown resample filter: " + name)
	}

	return
}

func ScaleName(filter *imaging.ResampleFilter) (name string) {
	switch filter {
	case &imaging.NearestNeighbor:
		name = "NearestNeighbor"
	case &imaging.Box:
		name = "Box"
	case &imaging.Linear:
		name = "Linear"
	case &imaging.Hermite:
		name = "Hermite"
	case &imaging.MitchellNetravali:
		name = "MitchellNetravali"
	case &imaging.CatmullRom:
		name = "CatmullRom"
	case &imaging.BSpline:
		name = "BSpline"
	case &imaging.Gaussian:
		name = "Gaussian"
	case &imaging.Bartlett:
		name = "Bartlett"
	case &imaging.Lanczos:
		name = "Lanczos"
	case &imaging.Hann:
		name = "Hann"
	case &imaging.Hamming:
		name = "Hamming"
	case &imaging.Blackman:
		name = "Blackman"
	case &imaging.Welch:
		name = "Welch"
	case &imaging.Cosine:
		name = "Cosine"
	default:
		panic("Unknown resample filter")
	}

	return
}
