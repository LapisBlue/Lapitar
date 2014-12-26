package mc

import "image"

type newSkin struct {
	image cuttableImage
	alex  bool
}

func (skin *newSkin) Image() image.Image {
	return skin.image
}

func (skin *newSkin) IsAlex() bool {
	return skin.alex
}

func (skin *newSkin) IsLegacy() bool {
	return false
}

func (skin *newSkin) Get(part SkinPart) image.Image {
	return skin.GetFace(part, All)
}

func (skin *newSkin) GetFace(part SkinPart, face Face) image.Image {
	return get(skin, skin.image, part, face)
}

func (skin *newSkin) Overlay(part SkinPart) image.Image {
	return skin.OverlayFace(part, All)
}

func (skin *newSkin) OverlayFace(part SkinPart, face Face) image.Image {
	return getOverlay(skin, skin.image, part, face)
}
