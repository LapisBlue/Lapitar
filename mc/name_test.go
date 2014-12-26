package mc

import "testing"

func TestNames(t *testing.T) {
	Steve()

	assert(t, IsName("Minecrell"))
	assert(t, IsName("Pande"))
	assert(t, IsName("Aesen_"))
	assert(t, IsName("A_b_C123698"))
	assert(t, IsName("A"))
	assert(t, IsName("BB"))
	assert(t, IsName("0123456789abcdef"))

	assert(t, !IsName("Invalid name"))
	assert(t, !IsName(".g lapitar"))
	assert(t, !IsName("*******"))
	assert(t, !IsName("thisisnotavalidnamebecauseitistoolong"))
	assert(t, !IsName(".,!"))
	assert(t, !IsName("Aesen!"))
}

func assert(t *testing.T, result bool) {
	if !result {
		t.Fail()
	}
}
