package mc

import "image"

type legacySkin struct {
	image cuttableImage
	alex  bool
}

func (skin *legacySkin) Image() image.Image {
	return skin.image
}

func (skin *legacySkin) IsAlex() bool {
	return skin.alex
}

func (skin *legacySkin) IsLegacy() bool {
	return true
}

func fallbackPart(part SkinPart) SkinPart {
	switch part {
	case LeftArm:
		return RightArm
	case LeftLeg:
		return RightLeg
	default:
		return part
	}
}

func (skin *legacySkin) Get(part SkinPart) image.Image {
	return skin.GetFace(part, All)
}

func (skin *legacySkin) GetFace(part SkinPart, face Face) image.Image {
	// TODO: flip arms and legs?
	return get(skin, skin.image, fallbackPart(part), face)
}

func (skin *legacySkin) Overlay(part SkinPart) image.Image {
	return skin.OverlayFace(part, All)
}

func (skin *legacySkin) OverlayFace(part SkinPart, face Face) image.Image {
	if part == Head {
		return getOverlay(skin, skin.image, part, face)
	} else {
		return nil
	}
}
