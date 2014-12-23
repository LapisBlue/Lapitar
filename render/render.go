package render

// #cgo pkg-config: gl glu osmesa
// #include "render.h"
import "C"
import (
	"errors"
	"github.com/LapisBlue/Lapitar/mc"
	"github.com/LapisBlue/Lapitar/util"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"unsafe"
)

const MinimumSize = 32

var DefaultScale = &imaging.Linear

var empty = make(map[image.Rectangle]*image.RGBA, 4)

func init() {
	empty[mc.HeadModel[0]] = image.NewRGBA(mc.HeadModel[0])
	empty[mc.BodyModel[0]] = image.NewRGBA(mc.BodyModel[0])
	empty[mc.LimbModel[0]] = image.NewRGBA(mc.LimbModel[0])
	empty[mc.AlexArmModel[0]] = image.NewRGBA(mc.AlexArmModel[0])
}

func Render(
	sk mc.Skin,
	angle float32, tilt float32, zoom float32,
	width, height int,
	superSampling int,
	portrait, full bool,
	overlay bool,
	shadow, lighting bool,
	filter *imaging.ResampleFilter) (result image.Image, err error) {

	w, h := width, height
	if superSampling > 1 {
		w = width * superSampling
		h = height * superSampling
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	head := prepare(sk, mc.Head)
	var headOverlay *image.RGBA
	if overlay {
		headOverlay = prepareOverlay(sk, mc.Head)
	}

	var body, bodyOverlay *image.RGBA
	var leftArm, leftArmOverlay, rightArm, rightArmOverlay *image.RGBA
	var leftLeg, leftLegOverlay, rightLeg, rightLegOverlay *image.RGBA

	if portrait || full {
		body = prepare(sk, mc.Body)
		leftArm = prepare(sk, mc.LeftArm)
		rightArm = prepare(sk, mc.RightArm)
		leftLeg = prepare(sk, mc.LeftLeg)
		rightLeg = prepare(sk, mc.RightLeg)

		if overlay {
			bodyOverlay = prepareOverlay(sk, mc.Body)
			leftArmOverlay = prepareOverlay(sk, mc.LeftArm)
			rightArmOverlay = prepareOverlay(sk, mc.RightArm)
			leftLegOverlay = prepareOverlay(sk, mc.LeftLeg)
			rightLegOverlay = prepareOverlay(sk, mc.RightLeg)
		}
	}

	// Render the head
	if !render(
		angle, tilt, zoom,
		shadow, lighting,
		portrait, full,
		overlay, !sk.IsLegacy(), sk.IsAlex(),
		img,
		head, headOverlay,
		body, bodyOverlay,
		leftArm, leftArmOverlay, rightArm, rightArmOverlay,
		leftLeg, leftLegOverlay, rightLeg, rightLegOverlay) {

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

func prepare(sk mc.Skin, part mc.SkinPart) *image.RGBA {
	return prepareUpload(sk.Get(part))
}

func prepareOverlay(sk mc.Skin, part mc.SkinPart) *image.RGBA {
	image := sk.Overlay(part)
	if image == nil || util.IsSolidColor(image) {
		return empty[mc.Model(sk, part)[0]]
	} else {
		return prepareUpload(image)
	}
}

type flippedImage struct {
	image image.Image
}

func (f *flippedImage) Bounds() image.Rectangle {
	return f.image.Bounds()
}

func (f *flippedImage) ColorModel() color.Model {
	return f.image.ColorModel()
}

func (f *flippedImage) At(x, y int) color.Color {
	return f.image.At(x, f.Bounds().Max.Y-y-1)
}

func prepareUpload(img image.Image) *image.RGBA {
	// This image is already fine for uploading
	if rgba, ok := img.(*image.RGBA); ok && rgba.Stride == rgba.Bounds().Dx()*4 {
		return rgba
	}

	// Convert image to RGBA
	rgba := image.NewRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba
}

func render(
	angle, tilt, zoom float32,
	shadow, lighting bool,
	portrait, full bool,
	overlay, newSkin, alex bool,
	result *image.RGBA,
	head, headOverlay *image.RGBA,
	body, bodyOverlay *image.RGBA,
	leftArm, leftArmOverlay, rightArm, rightArmOverlay *image.RGBA,
	leftLeg, leftLegOverlay, rightLeg, rightLegOverlay *image.RGBA) bool {

	return bool(
		C.Render(
			C.float(angle), C.float(tilt), C.float(zoom),
			C.bool(shadow), C.bool(lighting),
			C.bool(portrait), C.bool(full),
			C.bool(overlay), C.bool(newSkin), C.bool(alex),
			v_image(result),
			v_image(head), p_image(headOverlay),
			p_image(body), p_image(bodyOverlay),
			p_image(leftArm), p_image(leftArmOverlay), p_image(rightArm), p_image(rightArmOverlay),
			p_image(leftLeg), p_image(leftLegOverlay), p_image(rightLeg), p_image(rightLegOverlay)))
}

func v_image(image *image.RGBA) C.Image {
	return C.Image{unsafe.Pointer(&image.Pix[0]), C.int(image.Bounds().Dx()), C.int(image.Bounds().Dy())}
}

func p_image(image *image.RGBA) *C.Image {
	if image != nil {
		return &C.Image{unsafe.Pointer(&image.Pix[0]), C.int(image.Bounds().Dx()), C.int(image.Bounds().Dy())}
	} else {
		return nil
	}
}
