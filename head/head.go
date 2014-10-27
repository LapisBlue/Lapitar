package head

import (
	"github.com/LapisBlue/Tar/head/glctx"
	"github.com/LapisBlue/Tar/head/glu"
	"github.com/LapisBlue/Tar/skin"
	"github.com/go-gl/glow/gl/1.1/gl"
	"image"
	"unsafe"
)

type Renderer struct {
	GL glctx.ContextFactory

	Angle         float32
	Width, Height int
	w, h          int
	SuperSampling int

	Helmet           bool
	Shadow, Lighting bool
}

func (r *Renderer) Render(sk *skin.Skin) (head image.Image, err error) {
	if r.w == 0 || r.h == 0 {
		r.w = r.Width  //* r.SuperSampling
		r.h = r.Height //* r.SuperSampling
	}

	ctx, err := r.GL.Create(r.w, r.h)
	if err != nil {
		return
	}

	defer ctx.Close()

	gl.Init()

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1.0)
	gl.ShadeModel(gl.SMOOTH)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	glu.Perspective(45, float64(r.w)/float64(r.h), 0.1, 100)
	gl.MatrixMode(gl.MODELVIEW)

	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Render the head
	uploadImage(sk.Head(skin.All), 1)
	if r.Helmet {
		uploadImage(sk.Helm(skin.All), 2)
	}

	if r.Shadow {
		gl.Enable(gl.BLEND)
		gl.Disable(gl.TEXTURE_2D)
		gl.PushMatrix()
		gl.Translatef(0, -0.95, -0.45)

		count := float32(10)
		for i := float32(0); i < count; i++ {
			gl.Translatef(0, -0.01, 0)
			gl.Color4f(0, 0, 0, (1-(i/count))/2)
			r.draw(1.02, 0.01, 1.02)
		}

		gl.PopMatrix()
	}

	gl.Enable(gl.TEXTURE_2D)
	/*if r.Lighting { TODO
		gl.Enable(gl.LIGHTING)
		gl.Enable(gl.LIGHT0)
	}*/

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Color3f(1, 1, 1)
	gl.BindTexture(gl.TEXTURE_2D, 1)
	r.draw(1.05, 1.05, 1.05)

	if r.Helmet {
		gl.BindTexture(gl.TEXTURE_2D, 2)
		r.draw(1.05, 1.05, 1.05)
	}

	// TODO: Super sampling
	head = ctx.Render()
	return
}

func prepareUpload(img *image.RGBA) *image.RGBA {
	if img.Stride == img.Bounds().Dx()*4 {
		return img
	}

	// Convert image to RGBA
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

func uploadImage(img *image.RGBA, id uint32) {
	img = prepareUpload(img)
	gl.BindTexture(gl.TEXTURE_2D, id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), 0, gl.RGBA,
		gl.UNSIGNED_BYTE, unsafe.Pointer(&img.Pix[0]))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
}

func (r *Renderer) draw(x, y, z float32) {
	gl.PushMatrix()
	gl.Rotatef(20, 1, 0, 0)
	gl.Translatef(0, -1.5, -4.5)
	gl.Rotatef(r.Angle, 0, 1, 0)
	// TODO: Light
	// GL11.glLight(GL11.GL_LIGHT0, GL11.GL_POSITION, lightPosition);
	// GL11.glLight(GL11.GL_LIGHT0, GL11.GL_AMBIENT, lightAmbient);
	gl.Begin(gl.QUADS)
	gl.Normal3f(0, 0, -1)

	// Front
	gl.TexCoord2f(0.25, 1)
	gl.Vertex3f(-x, -y, z)
	gl.TexCoord2f(0.5, 1)
	gl.Vertex3f(x, -y, z)
	gl.TexCoord2f(0.5, 0.5)
	gl.Vertex3f(x, y, z)
	gl.TexCoord2f(0.25, 0.5)
	gl.Vertex3f(-x, y, z)

	// Back
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-x, -y, -z)
	gl.TexCoord2f(1, 0.5)
	gl.Vertex3f(-x, y, -z)
	gl.TexCoord2f(0.75, 0.5)
	gl.Vertex3f(x, y, -z)
	gl.TexCoord2f(0.75, 1)
	gl.Vertex3f(x, -y, -z)

	// Top
	gl.TexCoord2f(0.5, 0)
	gl.Vertex3f(-x, y, -z)
	gl.TexCoord2f(0.5, 0.5)
	gl.Vertex3f(-x, y, z)
	gl.TexCoord2f(0.25, 0.5)
	gl.Vertex3f(x, y, z)
	gl.TexCoord2f(0.25, 0)
	gl.Vertex3f(x, y, -z)

	// Bottom
	gl.TexCoord2f(0.5, 0.5)
	gl.Vertex3f(-x, -y, -z)
	gl.TexCoord2f(0.75, 0.5)
	gl.Vertex3f(x, -y, -z)
	gl.TexCoord2f(0.75, 0)
	gl.Vertex3f(x, -y, z)
	gl.TexCoord2f(0.5, 0)
	gl.Vertex3f(-x, -y, z)

	// Left
	gl.TexCoord2f(0.75, 1)
	gl.Vertex3f(x, -y, -z)
	gl.TexCoord2f(0.75, 0.5)
	gl.Vertex3f(x, y, -z)
	gl.TexCoord2f(0.5, 0.5)
	gl.Vertex3f(x, y, z)
	gl.TexCoord2f(0.5, 1)
	gl.Vertex3f(x, -y, z)

	// Right
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-x, -y, -z)
	gl.TexCoord2f(0.25, 1)
	gl.Vertex3f(-x, -y, z)
	gl.TexCoord2f(0.25, 0.5)
	gl.Vertex3f(-x, y, z)
	gl.TexCoord2f(0, 0.5)
	gl.Vertex3f(-x, y, -z)

	gl.End()
	gl.PopMatrix()
}
