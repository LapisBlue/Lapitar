package render

import (
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/disintegration/imaging"
	"image"
)

func RenderPortrait(
	sk mc.Skin,
	angle float32, tilt float32, zoom float32,
	width, height int,
	superSampling int,
	overlay bool,
	shadow, lighting bool,
	filter *imaging.ResampleFilter) (result image.Image, err error) {

	return
}
