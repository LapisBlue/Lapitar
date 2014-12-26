package mc

type Profile interface {
	Name() string
	UUID() string
	IsAlex() bool
}
