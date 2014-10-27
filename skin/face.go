package skin

type Face byte

const (
	All = Face(iota)
	Top
	Bottom
	Right
	Front
	Left
	Back
	faceCount
)
