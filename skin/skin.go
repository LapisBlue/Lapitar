package skin

import (
	"image"
)

type Skin image.RGBA

func (skin *Skin) Image() *image.RGBA {
	return (*image.RGBA)(skin)
}

// Head
var head = [faceCount]image.Rectangle{
	rectFrom(0, 0, 32, 16), // All
	rectFrom(8, 0, 8, 8),   // Top
	rectFrom(16, 0, 8, 8),  // Bottom
	rectFrom(0, 8, 8, 8),   // Right
	rectFrom(8, 8, 8, 8),   // Front
	rectFrom(16, 8, 8, 8),  // Left
	rectFrom(24, 8, 8, 8),  // Back
}

func (skin *Skin) Head(face Face) *image.RGBA {
	return rgba(skin.Image().SubImage(head[face]))
}

// Helm
var helm = [faceCount]image.Rectangle{
	rectFrom(32, 0, 32, 16), // All
	rectFrom(40, 0, 8, 8),   // Top
	rectFrom(48, 0, 8, 8),   // Bottom
	rectFrom(32, 8, 8, 8),   // Right
	rectFrom(40, 8, 8, 8),   // Front
	rectFrom(48, 8, 8, 8),   // Left
	rectFrom(56, 8, 8, 8),   // Back
}

func (skin *Skin) Helm(face Face) *image.RGBA {
	return rgba(skin.Image().SubImage(helm[face]))
}

// Body
var body = [faceCount]image.Rectangle{
	rectFrom(16, 16, 24, 16), // All
	rectFrom(20, 16, 8, 4),   // Top
	rectFrom(28, 16, 8, 4),   // Bottom
	rectFrom(16, 20, 4, 12),  // Right
	rectFrom(20, 20, 8, 12),  // Front
	rectFrom(28, 20, 4, 12),  // Left
	rectFrom(32, 20, 8, 12),  // Back
}

func (skin *Skin) Body(face Face) *image.RGBA {
	return rgba(skin.Image().SubImage(body[face]))
}

// Arm
var arm = [faceCount]image.Rectangle{
	rectFrom(40, 16, 16, 16), // All
	rectFrom(44, 16, 4, 4),   // Top
	rectFrom(48, 16, 4, 4),   // Bottom
	rectFrom(40, 20, 4, 12),  // Right
	rectFrom(44, 20, 4, 12),  // Front
	rectFrom(48, 20, 4, 12),  // Left
	rectFrom(52, 20, 4, 12),  // Back
}

func (skin *Skin) Arm(face Face) *image.RGBA {
	return rgba(skin.Image().SubImage(arm[face]))
}

// Leg
var leg = [faceCount]image.Rectangle{
	rectFrom(0, 16, 16, 16), // All
	rectFrom(4, 16, 4, 4),   // Top
	rectFrom(8, 16, 4, 4),   // Bottom
	rectFrom(0, 20, 4, 12),  // Right
	rectFrom(4, 20, 4, 12),  // Front
	rectFrom(8, 20, 4, 12),  // Left
	rectFrom(12, 20, 4, 12), // Back
}

func (skin *Skin) Leg(face Face) *image.RGBA {
	return rgba(skin.Image().SubImage(leg[face]))
}

func rectFrom(x, y, width, height int) image.Rectangle {
	return image.Rect(x, y, x+width, y+height)
}

func rgba(img image.Image) (rgba *image.RGBA) {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	// Convert image to RGBA
	rgba = image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return
}
