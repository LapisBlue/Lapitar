package glctx

import (
	"image"
	"io"
)

type Context interface {
	Size() image.Point
	Render() image.Image
	io.Closer
}

type ContextFactory interface {
	Create(width, height int) (Context, error)
}
