package cache

import "github.com/LapisBlue/Lapitar/mc"

func Fallback(profile mc.Profile) SkinMeta {
	var handle *defaultSkinMeta
	if profile.IsAlex() {
		handle = alex
	} else {
		handle = steve
	}

	return &fallbackSkinMeta{handle, profile}
}

type fallbackSkinMeta struct {
	*defaultSkinMeta
	profile mc.Profile
}

func (meta *fallbackSkinMeta) Profile() mc.Profile {
	return meta.profile
}

func FallbackByUUID(uuid string) SkinMeta {
	return Fallback(uuidProfile(uuid))
}

type uuidProfile string

func (uuid uuidProfile) Name() string {
	return "Unknown (" + uuid.UUID() + ")"
}

func (uuid uuidProfile) UUID() string {
	return string(uuid)
}

func (uuid uuidProfile) IsAlex() bool {
	return mc.IsAlex(uuid.UUID())
}

type nameProfile string

func (name nameProfile) Name() string {
	return string(name)
}

func (name nameProfile) UUID() string {
	panic("Unknown UUID")
}

func (uuid nameProfile) IsAlex() bool {
	return false // Use steve by default
}

func FallbackByName(name string) SkinMeta {
	return Fallback(nameProfile(name))
}
