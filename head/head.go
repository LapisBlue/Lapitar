package head

// #cgo pkg-config: gl glu osmesa
// #include "head.h"
import "C"
import (
	"errors"
	"github.com/LapisBlue/Tar/skin"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"unsafe"
)

func Render(
	sk *skin.Skin,
	angle float32,
	width, height int,
	superSampling int,
	helmet bool,
	shadow, lighting bool) (result image.Image, err error) {

	w, h := width, height
	if superSampling > 1 {
		w = width * superSampling
		h = height * superSampling
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	head := prepareUpload(sk.Head(skin.All))
	var helm *image.RGBA
	if helmet {
		helm = sk.Helm(skin.All)
		if isSolidColor(helm) {
			helm = nil
			helmet = false
		} else {
			helm = prepareUpload(sk.Helm(skin.All))
		}
	}

	// Render the head
	if !renderHead(angle, shadow, lighting, img, head, helm) {
		return nil, errors.New("Rendering failed")
	}

	result = img
	if superSampling > 1 {
		result = resize.Resize(uint(width), uint(height), result, resize.Bicubic)
	}

	// The result is flipped tbh
	result = &flippedImage{result}
	return
}

func renderHead(
	angle float32,
	shadow, lighting bool,
	result *image.RGBA,
	head *image.RGBA,
	helm *image.RGBA) bool {

	var helmPointer unsafe.Pointer
	helmWidth, helmHeight := 0, 0
	if helm != nil {
		helmPointer = unsafe.Pointer(&helm.Pix[0])
		helmWidth = helm.Bounds().Dx()
		helmHeight = helm.Bounds().Dy()
	}

	return bool(
		C.RenderHead(
			C.float(angle),
			C.bool(shadow), C.bool(lighting),
			unsafe.Pointer(&result.Pix[0]), C.int(result.Bounds().Dx()), C.int(result.Bounds().Dy()),
			unsafe.Pointer(&head.Pix[0]), C.int(head.Bounds().Dx()), C.int(head.Bounds().Dy()),
			helmPointer, C.int(helmWidth), C.int(helmHeight)))
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

func isSolidColor(img image.Image) bool {
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

func prepareUpload(img *image.RGBA) *image.RGBA {
	if img.Stride == img.Bounds().Dx()*4 {
		return img
	}

	// While this image view (e.g. through SubImage) is faster, we need exactly this image only in memory for OpenGL
	rgba := image.NewRGBA(img.Bounds())
	pos := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			from := img.PixOffset(x, y)
			for i := 0; i < 4; i, pos = i+1, pos+1 {
				rgba.Pix[pos] = img.Pix[from+i]
			}
		}
	}

	return rgba
}
