package main

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
)

// Version is a type for storing ILCD version information
type Version struct {
	Major  int
	Minor  int
	Update int
}

// ParseVersion parses a ILCD version string (e.g. 01.24.001)
func ParseVersion(text string) *Version {
	v := Version{}
	part := 0
	var buf bytes.Buffer
	addPart := func() {
		i := 0
		if buf.Len() > 0 {
			s := buf.String()
			i, _ = strconv.Atoi(s)
		}
		switch part {
		case 0:
			v.Major = i
		case 1:
			v.Minor = i
		case 2:
			v.Update = i
		}
		buf.Reset()
		part++
	}
	for _, r := range strings.TrimSpace(text) {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else if r == '.' {
			addPart()
		}
	}
	addPart()
	return &v
}

// String returns the version in the ILCD format (e.g. 01.00.000).
func (v *Version) String() string {
	if v == nil {
		return "00.00.000"
	}
	return pad(v.Major, 2) + "." + pad(v.Minor, 2) + "." + pad(v.Update, 3)
}

// NormalizeVersion converts the given version into the ILCD format, e.g. `1`
// and `1.0` become `01.00.000`. A blank version stays blank.
func NormalizeVersion(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	return ParseVersion(text).String()
}

// pad converts the given number into a string with leading zeros.
func pad(value int, length int) string {
	s := strconv.Itoa(value)
	for len(s) < length {
		s = "0" + s
	}
	return s
}

// Compare compares the version with another version.
func (v *Version) Compare(other *Version) int {
	if other == nil {
		return 1
	}
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Update - other.Update
}

// Same returns true if both versions are exactly the same.
func (v *Version) Same(other *Version) bool {
	return v.Compare(other) == 0
}

// NewerThan returns true if the version is new (mean higher) than the given
// version.
func (v *Version) NewerThan(other *Version) bool {
	return v.Compare(other) > 0
}
