package glctx

import (
	"github.com/LapisBlue/Tar/head/glctx/osmesa"
	"image"
	"unsafe"
)

type osMesaError string

func (err osMesaError) Error() string {
	return "osmesa: operation failed: " + string(err)
}

type osMesaContext struct {
	ctx    osmesa.Context
	render *image.RGBA
	size   image.Point
}

func (mesa *osMesaContext) Size() image.Point {
	return mesa.size
}

func (mesa *osMesaContext) Render() image.Image {
	return mesa.render
}

func (mesa *osMesaContext) Close() error {
	osmesa.DestroyContext(mesa.ctx)
	return nil
}

type osMesaContextFactory struct{}

func (factory *osMesaContextFactory) Create(width, height int) (Context, error) {
	mesa := &osMesaContext{
		size: image.Pt(width, height),
	}

	mesa.ctx = osmesa.CreateContextExt(osmesa.RGBA, 16, 0, 0, nil)
	if mesa.ctx == nil {
		return nil, osMesaError("CreateContextExt")
	}

	mesa.render = image.NewRGBA(image.Rect(0, 0, width, height))
	if !osmesa.MakeCurrent(mesa.ctx, unsafe.Pointer(&mesa.render.Pix[0]), osmesa.UNSIGNED_BYTE, width, height) {
		return nil, osMesaError("MakeCurrent")
	}

	return mesa, nil
}

var osMesa = &osMesaContextFactory{}

func OSMesa() ContextFactory {
	return osMesa
}
