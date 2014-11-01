package face

import (
	"github.com/LapisBlue/Lapitar/skin"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
)

func Render(sk *skin.Skin, size int, helmet bool) image.Image {
	face := sk.Head(skin.Front)
	if helmet {
		helm := sk.Helm(skin.Front)
		if !util.IsSolidColor(helm) {
			temp := imaging.Clone(face)
			draw.Draw(temp, face.Bounds(), helm, helm.Bounds().Min, draw.Over)
			face = temp
		}
	}

	return imaging.Resize(face, size, size, imaging.NearestNeighbor)
}
