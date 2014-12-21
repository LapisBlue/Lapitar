package render

import (
	"errors"
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
	"image"
)

func RenderHead(
	sk mc.Skin,
	angle float32, tilt float32, zoom float32,
	width, height int,
	superSampling int,
	overlay bool,
	shadow, lighting bool,
	filter *imaging.ResampleFilter) (result image.Image, err error) {

	w, h := width, height
	if superSampling > 1 {
		w = width * superSampling
		h = height * superSampling
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	head := prepareUpload(sk.Get(mc.Head))
	var helm *image.RGBA
	if overlay {
		helmImg := sk.Overlay(mc.Head)
		if util.IsSolidColor(helmImg) {
			helm = nil
			overlay = false
		} else {
			helm = prepareUpload(helmImg)
		}
	}

	// Render the head
	if !render(
		angle, tilt, zoom,
		shadow, lighting,
		false, false,
		overlay,
		!sk.IsLegacy(),
		img,
		head, helm,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil) {

		return nil, errors.New("Rendering failed")
	}

	result = img
	if superSampling > 1 {
		result = imaging.Resize(result, width, height, *filter)
	}

	// The result is flipped tbh
	result = &flippedImage{result}
	return
}
