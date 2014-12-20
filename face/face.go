package face

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
)

const MinimumSize = 8

var DefaultScale = &imaging.NearestNeighbor

func Render(sk *mc.Skin, size int, helmet bool, filter *imaging.ResampleFilter) image.Image {
	face := sk.Head(mc.Front)
	if helmet {
		helm := sk.Helm(mc.Front)
		if !util.IsSolidColor(helm) {
			temp := imaging.Clone(face)
			draw.Draw(temp, face.Bounds(), helm, helm.Bounds().Min, draw.Over)
			face = temp
		}
	}

	if size == MinimumSize {
		return face
	}

	return imaging.Resize(face, size, size, *filter)
}
