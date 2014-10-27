/*
 * Mesa 3-D graphics library
 * Version:  6.5
 *
 * Copyright (C) 1999-2005  Brian Paul   All Rights Reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
 * OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * BRIAN PAUL BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN
 * AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

/*
 * Mesa Off-Screen rendering interface.
 *
 * This is an operating system and window system independent interface to
 * Mesa which allows one to render images into a client-supplied buffer in
 * main memory.  Such images may manipulated or saved in whatever way the
 * client wants.
 *
 * These are the API functions:
 *   osmesa.CreateContext - create a new Off-Screen Mesa rendering context
 *   osmesa.MakeCurrent - bind an osmesa.Context to a client's image buffer
 *                       and make the specified context the current one.
 *   osmesa.DestroyContext - destroy an osmesa.Context
 *   osmesa.GetCurrentContext - return thread's current context ID
 *   osmesa.PixelStore - controls how pixels are stored in image buffer
 *   osmesa.GetIntegerv - return OSMesa state parameters
 *
 *
 * The limits on the width and height of an image buffer are MAX_WIDTH and
 * MAX_HEIGHT as defined in Mesa/src/config.h.  Defaults are 1280 and 1024.
 * You can increase them as needed but beware that many temporary arrays in
 * Mesa are dimensioned by MAX_WIDTH or MAX_HEIGHT.
 */

package osmesa

// #cgo pkg-config: osmesa
// #include <GL/osmesa.h>
import "C"
import "github.com/go-gl/glow/gl/1.1/gl"
import "unsafe"

type PixelType int

const (
	UNSIGNED_BYTE = gl.UNSIGNED_BYTE
)

type Format int

/*
 * Values for the format parameter of OSMesaCreateContext()
 * New in version 2.0.
 */
const (
	COLOR_INDEX Format = gl.COLOR_INDEX
	RGBA        Format = gl.RGBA
	BGRA        Format = 0x1
	ARGB        Format = 0x2
	RGB         Format = gl.RGB
	BGR         Format = 0x4
	RGB_565     Format = 0x5
)

type PixelStoreParameter int

/*
 * OSMesaPixelStore() parameters:
 * New in version 2.0.
 */
const (
	ROW_LENGTH PixelStoreParameter = 0x10
	Y_UP       PixelStoreParameter = 0x11
)

/*
 * Accepted by OSMesaGetIntegerv:
 */
type PropertyName int

const (
	WIDTH      PropertyName = 0x20
	HEIGHT     PropertyName = 0x21
	FORMAT     PropertyName = 0x22
	TYPE       PropertyName = 0x23
	MAX_WIDTH  PropertyName = 0x24 /* new in 4.0 */
	MAX_HEIGHT PropertyName = 0x25 /* new in 4.0 */
)

type Context C.OSMesaContext

/*
 * Create an Off-Screen Mesa rendering context.  The only attribute needed is
 * an RGBA vs Color-Index mode flag.
 *
 * Input:  format - one of OSMESA_COLOR_INDEX, OSMESA_RGBA, OSMESA_BGRA,
 *                  OSMESA_ARGB, OSMESA_RGB, or OSMESA_BGR.
 *         sharelist - specifies another osmesa.Context with which to share
 *                     display lists.  NULL indicates no sharing.
 * Return:  an osmesa.Context or 0 if error
 */
func CreateContext(format Format, sharelist Context) Context {
	return Context(C.OSMesaCreateContext(C.GLenum(format), C.OSMesaContext(sharelist)))

}

/*
 * Create an Off-Screen Mesa rendering context and specify desired
 * size of depth buffer, stencil buffer and accumulation buffer.
 * If you specify zero for depthBits, stencilBits, accumBits you
 * can save some memory.
 *
 * New in Mesa 3.5
 */
func CreateContextExt(format Format, depthBits int, stencilBits int,
	accumBits int, sharelist Context) Context {
	return Context(C.OSMesaCreateContextExt(C.GLenum(format), C.GLint(depthBits), C.GLint(stencilBits),
		C.GLint(accumBits), C.OSMesaContext(sharelist)))
}

/*
 * Destroy an Off-Screen Mesa rendering context.
 *
 * Input:  ctx - the context to destroy
 */
func DestroyContext(ctx Context) {
	C.OSMesaDestroyContext(C.OSMesaContext(ctx))
}

/*
 * Bind an osmesa.Context to an image buffer.  The image buffer is just a
 * block of memory which the client provides.  Its size must be at least
 * as large as width*height*sizeof(type).  Its address should be a multiple
 * of 4 if using RGBA mode.
 *
 * Image data is stored in the order of glDrawPixels:  row-major order
 * with the lower-left image pixel stored in the first array position
 * (ie. bottom-to-top).
 *
 * Since the only type initially supported is GL_UNSIGNED_BYTE, if the
 * context is in RGBA mode, each pixel will be stored as a 4-byte RGBA
 * value.  If the context is in color indexed mode, each pixel will be
 * stored as a 1-byte value.
 *
 * If the context's viewport hasn't been initialized yet, it will now be
 * initialized to (0,0,width,height).
 *
 * Input:  ctx - the rendering context
 *         buffer - the image buffer memory
 *         type - data type for pixel components, only GL_UNSIGNED_BYTE
 *                supported now
 *         width, height - size of image buffer in pixels, at least 1
 * Return:  GL_TRUE if success, GL_FALSE if error because of invalid ctx,
 *          invalid buffer address, type!=GL_UNSIGNED_BYTE, width<1, height<1,
 *          width>internal limit or height>internal limit.
 */
func MakeCurrent(ctx Context, buffer unsafe.Pointer, type_ PixelType,
	width int, height int) bool {

	return 0 != C.OSMesaMakeCurrent(C.OSMesaContext(ctx), buffer, C.GLenum(type_),
		C.GLsizei(width), C.GLsizei(height))
}

/*
 * Return the current Off-Screen Mesa rendering context handle.
 */
func GetCurrentContext() Context {
	return Context(C.OSMesaGetCurrentContext())
}

/*
 * Set pixel store/packing parameters for the current context.
 * This is similar to glPixelStore.
 * Input:  pname - OSMESA_ROW_LENGTH
 *                    specify actual pixels per row in image buffer
 *                    0 = same as image width (default)
 *                 OSMESA_Y_UP
 *                    zero = Y coordinates increase downward
 *                    non-zero = Y coordinates increase upward (default)
 *         value - the value for the parameter pname
 *
 * New in version 2.0.
 */
func PixelStore(pname PixelStoreParameter, value int) {
	C.OSMesaPixelStore(C.GLint(pname), C.GLint(value))
}

/*
 * Return an integer value like glGetIntegerv.
 * Input:  pname -
 *                 OSMESA_WIDTH  return current image width
 *                 OSMESA_HEIGHT  return current image height
 *                 OSMESA_FORMAT  return image format
 *                 OSMESA_TYPE  return color component data type
 *                 OSMESA_ROW_LENGTH return row length in pixels
 *                 OSMESA_Y_UP returns 1 or 0 to indicate Y axis direction
 *         value - pointer to integer in which to return result.
 */
func GetIntegerv(pname PropertyName) int {
	var value C.GLint
	C.OSMesaGetIntegerv(C.GLint(pname), (*C.GLint)(&value))
	return int(value)
}

/*
 * Return the depth buffer associated with an OSMesa context.
 * Input:  c - the OSMesa context
 * Output:  width, height - size of buffer in pixels
 *          bytesPerValue - bytes per depth value (2 or 4)
 *          buffer - pointer to depth buffer values
 * Return:  GL_TRUE or GL_FALSE to indicate success or failure.
 *
 * New in Mesa 2.4.
 */
func GetDepthBuffer(c Context, width *int, height *int,
	bytesPerValue *int, buffer *unsafe.Pointer) {

	var width_ C.GLint
	var height_ C.GLint
	var bytesPerValue_ C.GLint
	var buffer_ unsafe.Pointer

	C.OSMesaGetDepthBuffer(C.OSMesaContext(c), &width_, &height_,
		&bytesPerValue_, &buffer_)

	*width = int(width_)
	*height = int(height_)
	*bytesPerValue = int(bytesPerValue_)
	*buffer = buffer_

}

/*
 * Return the color buffer associated with an OSMesa context.
 * Input:  c - the OSMesa context
 * Output:  width, height - size of buffer in pixels
 *          format - buffer format (OSMESA_FORMAT)
 *          buffer - pointer to depth buffer values
 * Return:  GL_TRUE or GL_FALSE to indicate success or failure.
 *
 * New in Mesa 3.3.
 */
func GetColorBuffer(c Context, width *int, height *int,
	format *Format, buffer *unsafe.Pointer) {

	var width_ C.GLint
	var height_ C.GLint
	var format_ C.GLint
	var buffer_ unsafe.Pointer

	C.OSMesaGetDepthBuffer(C.OSMesaContext(c), &width_, &height_,
		&format_, &buffer_)

	*width = int(width_)
	*height = int(height_)
	*format = Format(format_)
	*buffer = buffer_

}

/**
 * Enable/disable color clamping, off by default.
 * New in Mesa 6.4.2
 */
func ColorClamp(enable bool) {
	C.OSMesaColorClamp(glBool(enable))
}

func glBool(v bool) C.GLboolean {
	if v {
		return 1
	}
	return 0
}
