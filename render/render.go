package render

// #cgo pkg-config: gl glu osmesa
// #include "render.h"
import "C"
import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"unsafe"
)

const MinimalSize = 32

var DefaultScale = &imaging.Linear

func render(
	angle, tilt, zoom float32,
	shadow, lighting bool,
	portrait, full bool,
	overlay, newSkin bool,
	result *image.RGBA,
	head *image.RGBA, headOverlay *image.RGBA,
	body *image.RGBA, bodyOverlay *image.RGBA,
	leftArm *image.RGBA, leftArmOverlay *image.RGBA, rightArm *image.RGBA, rightArmOverlay *image.RGBA,
	leftLeg *image.RGBA, leftLegOverlay *image.RGBA, rightLeg *image.RGBA, rightLegOverlay *image.RGBA) bool {

	return bool(
		C.Render(
			C.float(angle), C.float(tilt), C.float(zoom),
			C.bool(shadow), C.bool(lighting),
			C.bool(portrait), C.bool(full),
			C.bool(overlay), C.bool(newSkin),
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
