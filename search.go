package main

import "strings"

// filterByName applies a simple name search to the given index entries. The
// search phrase is split into whitespace separated keywords and an entry
// matches when all keywords are contained in its name (case-insensitive). A
// phrase that contains the prefix `uuid:` is interpreted as a UUID prefix
// search instead. A blank phrase returns the entries unchanged.
func filterByName(entries []*indexEntry, phrase string) []*indexEntry {
	phrase = strings.TrimSpace(phrase)
	if phrase == "" {
		return entries
	}

	// a phrase with a `uuid:` prefix searches for the UUID instead
	lower := strings.ToLower(phrase)
	if i := strings.Index(lower, "uuid:"); i >= 0 {
		prefix := searchUUID(lower[i+len("uuid:"):])
		if prefix != "" {
			filtered := make([]*indexEntry, 0, len(entries))
			for _, e := range entries {
				if strings.HasPrefix(searchUUID(e.UUID), prefix) {
					filtered = append(filtered, e)
				}
			}
			return filtered
		}
	}

	keywords := strings.Fields(phrase)
	if len(keywords) == 0 {
		return entries
	}
	for i := range keywords {
		keywords[i] = strings.ToLower(keywords[i])
	}

	filtered := make([]*indexEntry, 0, len(entries))
	for _, e := range entries {
		name := strings.ToLower(e.Name)
		match := true
		for _, keyword := range keywords {
			if !strings.Contains(name, keyword) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// searchUUID prepares a UUID or a UUID prefix for a comparison: separators are
// removed so that e.g. `2222-2222` also matches a UUID that starts with
// `22222222`.
func searchUUID(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	return strings.ReplaceAll(text, "-", "")
}
