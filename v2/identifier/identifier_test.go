// Tideland Go Application Support - Identifier - Unit Tests
//
// Copyright (C) 2009-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package identifier_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"github.com/tideland/goas/v2/identifier"
	"github.com/tideland/gots/v3/asserts"
)

//--------------------
// TESTS
//--------------------

// Test the UUID.
func TestUUID(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)
	// Asserts.
	uuid := identifier.NewUUID()
	uuidStr := uuid.String()
	assert.Equal(len(uuid), 16, "UUID length has to be 16")
	assert.Match(uuidStr, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", "UUID has to match")
	// Check for unique creation, but only weak test.
	uuids := make(map[string]bool)
	for i := 0; i < 1000000; i++ {
		uuid = identifier.NewUUID()
		uuidStr = uuid.String()
		assert.False(uuids[uuidStr], "UUID collision must not happen")
		uuids[uuidStr] = true
	}
	// Check for copy.
	uuidA := identifier.NewUUID()
	uuidB := uuidA.Copy()
	for i := 0; i < len(uuidA); i++ {
		uuidA[i] = 0
	}
	assert.Different(uuidA, uuidB)
}

// Test the creation of identifiers based on types.
func TestTypeAsIdentifierPart(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	// Type as identifier.
	var tai TypeToSplitForIdentifier

	id := identifier.TypeAsIdentifierPart(tai)
	assert.Equal(id, "type-to-split-for-identifier", "wrong TypeAsIdentifierPart() result")

	id = identifier.TypeAsIdentifierPart(identifier.NewUUID())
	assert.Equal(id, "u-u-i-d", "wrong TypeAsIdentifierPart() result")
}

// Test the creation of identifiers based on parts.
func TestIdentifier(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	// Identifier.
	id := identifier.Identifier("One", 2, "three four")
	assert.Equal(id, "one:2:three-four", "wrong Identifier() result")

	id = identifier.Identifier(2011, 6, 22, "One, two, or  three things.")
	assert.Equal(id, "2011:6:22:one-two-or-three-things", "wrong Identifier() result")
}

// Test the creation of identifiers based on parts with defined seperators.
func TestSepIdentifier(t *testing.T) {
	assert := asserts.NewTestingAssertion(t, true)

	id := identifier.SepIdentifier("+", 1, "oNe", 2, "TWO", "3", "ÄÖÜ")
	assert.Equal(id, "1+one+2+two+3+äöü", "wrong SepIdentifier() result")

	id = identifier.LimitedSepIdentifier("+", true, "     ", 1, "oNe", 2, "TWO", "3", "ÄÖÜ", "Four", "+#-:,")
	assert.Equal(id, "1+one+2+two+3+four", "wrong LimitedSepIdentifier() result")
}

//--------------------
// HELPER
//--------------------

// Type as part of an identifier.
type TypeToSplitForIdentifier bool

// EOF
