package mc

type SkinProfile interface {
	Profile
	Skin() SkinMeta
}

type SkinMeta interface {
	ID() string
	URL() string
	Download() (Skin, error)
}
