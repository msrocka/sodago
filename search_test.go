package main

import "testing"

func TestFilterByName(t *testing.T) {
	entries := []*indexEntry{
		{UUID: "11111111-1111-1111-1111-111111111111", Name: "ACME Corporation"},
		{UUID: "22222222-2222-2222-2222-222222222222", Name: "Beta Consulting"},
		{UUID: "33333333-3333-3333-3333-333333333333", Name: "acme labs"},
	}

	// a blank phrase returns all entries
	if got := filterByName(entries, ""); len(got) != 3 {
		t.Error("expected all entries for a blank phrase")
	}
	if got := filterByName(entries, "   "); len(got) != 3 {
		t.Error("expected all entries for a blank phrase")
	}

	check := func(phrase string, expected ...string) {
		filtered := filterByName(entries, phrase)
		if len(filtered) != len(expected) {
			t.Fatalf("filter %q: expected %d results but got %d",
				phrase, len(expected), len(filtered))
		}
		for i, e := range filtered {
			if e.UUID != expected[i] {
				t.Errorf("filter %q: expected %s at %d but got %s",
					phrase, expected[i], i, e.UUID)
			}
		}
	}

	// case-insensitive substring match
	check("acme", entries[0].UUID, entries[2].UUID)
	check("ACME", entries[0].UUID, entries[2].UUID)
	check("Corp", entries[0].UUID)

	// multiple keywords must all match
	check("acme corp", entries[0].UUID)
	check("acme consulting Beta")

	// UUID prefix search
	check("uuid:2222", entries[1].UUID)
	check("UUID:2222", entries[1].UUID)
	check("uuid:2222-2222", entries[1].UUID)

	// no match
	check("does not exist")
}
