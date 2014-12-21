package mc

import (
	"image"
)

type SkinPart byte

const (
	Head SkinPart = iota
	Body
	LeftArm
	RightArm
	LeftLeg
	RightLeg
	skinPartCount
)

type Face byte

const (
	All Face = iota
	Front
	Back
	Top
	Bottom
	Left
	Right
	faceCount
)

var HeadModel = [faceCount]image.Rectangle{
	bounds(32, 16),        // All
	relative(8, 8, 8, 8),  // Front
	relative(24, 8, 8, 8), // Back
	relative(8, 0, 8, 8),  // Top
	relative(16, 0, 8, 8), // Bottom
	relative(16, 8, 8, 8), // Left
	relative(0, 8, 8, 8),  // Right
}

var BodyModel = [faceCount]image.Rectangle{
	bounds(24, 16),         // All
	relative(4, 4, 8, 12),  // Front
	relative(16, 4, 8, 12), // Back
	relative(4, 0, 8, 4),   // Top
	relative(12, 0, 8, 4),  // Bottom
	relative(12, 4, 4, 12), // Left
	relative(0, 4, 4, 12),  // Right
}

var LimbModel = [faceCount]image.Rectangle{
	bounds(16, 16),         // All
	relative(4, 4, 4, 12),  // Front
	relative(12, 4, 4, 12), // Back
	relative(4, 0, 4, 4),   // Top
	relative(8, 0, 4, 4),   // Bottom
	relative(8, 4, 4, 12),  // Left
	relative(0, 4, 4, 12),  // Right
}

var AlexArmModel = [faceCount]image.Rectangle{
	bounds(14, 16),         // All
	relative(4, 4, 3, 12),  // Front
	relative(11, 4, 3, 12), // Back
	relative(4, 0, 3, 4),   // Top
	relative(7, 0, 3, 4),   // Bottom
	relative(7, 4, 4, 12),  // Left
	relative(0, 4, 4, 12),  // Right
}

var Positions = [skinPartCount]image.Point{
	image.ZP,         // Head
	image.Pt(16, 16), // Body
	image.Pt(32, 48), // Left Arm
	image.Pt(40, 16), // Right Arm
	image.Pt(16, 48), // Left Leg
	image.Pt(0, 16),  // Right Leg
}

var OverlayPositions = [skinPartCount]image.Point{
	image.Pt(32, 0),  // Head Overlay
	image.Pt(16, 32), // Body Overlay
	image.Pt(48, 48), // Left Arm Overlay
	image.Pt(40, 32), // Right Arm Overlay
	image.Pt(0, 48),  // Left Leg Overlay
	image.Pt(0, 32),  // Right Leg Overlay
}

type Skin interface {
	Image() image.Image
	IsAlex() bool
	IsLegacy() bool
	Get(part SkinPart) image.Image
	GetFace(part SkinPart, face Face) image.Image
	Overlay(part SkinPart) image.Image
	OverlayFace(part SkinPart, face Face) image.Image
}

type cuttableImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

func CreateSkin(skin image.Image, alex bool) Skin {
	if image, ok := skin.(cuttableImage); ok {
		height := image.Bounds().Dy()
		if height == 64 {
			return &newSkin{image: image, alex: alex}
		} else if height == 32 {
			return &legacySkin{image: image, alex: alex}
		} else {
			panic("Unsupported skin format")
		}
	}

	panic("Unsupported image type") // TODO
}

func Model(skin Skin, part SkinPart) [faceCount]image.Rectangle {
	switch part {
	case Head:
		return HeadModel
	case Body:
		return BodyModel
	case LeftLeg, RightLeg:
		return LimbModel
	case RightArm, LeftArm:
		if skin != nil && skin.IsAlex() {
			return AlexArmModel
		} else {
			return LimbModel
		}
	default:
		panic("Unknown skin part")
	}
}

func get(skin Skin, image cuttableImage, part SkinPart, face Face) image.Image {
	return image.SubImage(Model(skin, part)[face].Add(Positions[part]))
}

func getOverlay(skin Skin, image cuttableImage, part SkinPart, face Face) image.Image {
	return image.SubImage(Model(skin, part)[face].Add(OverlayPositions[part]))
}

func bounds(x, y int) image.Rectangle {
	return relative(0, 0, x, y)
}

func relative(x, y, width, height int) image.Rectangle {
	return image.Rect(x, y, x+width, y+height)
}
