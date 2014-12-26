package mc

import "testing"

func TestUUIDs(t *testing.T) {
	assert(t, IsUUID("393f9948d7a24d848617ee28840780f1"))
	assert(t, IsUUID("f66649bf733248f889b8d87cea51a745"))
	assert(t, IsUUID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	assert(t, IsUUID("99999999999999999999999999999999"))
	assert(t, IsUUID(ParseUUID("c4ddb997-f7a8-4257-9209-ab066b6c376c")))
	assert(t, IsUUID(ParseUUID("b43ecf24-f0a4-4bba-89a5-04fd6833eb92")))
	assert(t, IsUUID(ParseUUID("5-e-9-c-8-9-9-3-9-4-d-2-4-d-8-8-b-c-b-6-2-3-d-0-1-e-8-3-e-6-c-1")))
	assert(t, IsUUID(ParseUUID("393F9948D7A24D848617EE28840780F1")))

	assert(t, !IsUUID("393F9948D7A24D848617EE28840780F1"))
	assert(t, !IsUUID("Invalid UUID"))
	assert(t, !IsUUID(".g lapitar"))
	assert(t, !IsUUID("a_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))
}
