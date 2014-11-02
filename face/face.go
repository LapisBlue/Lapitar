package face

import (
	"github.com/LapisBlue/Lapitar/skin"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
)

const MinimalSize = 8

var DefaultScale = &imaging.NearestNeighbor

func Render(sk *skin.Skin, size int, helmet bool, filter *imaging.ResampleFilter) image.Image {
	face := sk.Head(skin.Front)
	if helmet {
		helm := sk.Helm(skin.Front)
		if !util.IsSolidColor(helm) {
			temp := imaging.Clone(face)
			draw.Draw(temp, face.Bounds(), helm, helm.Bounds().Min, draw.Over)
			face = temp
		}
	}

	if size == MinimalSize {
		return face
	}

	return imaging.Resize(face, size, size, *filter)
}
