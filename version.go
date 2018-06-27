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
			s := string(buf.Bytes())
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

func (v *Version) String() string {
	if v == nil {
		return "00.00.000"
	}
	major := strconv.Itoa(v.Major)
	if len(major) == 1 {
		major = "0" + major
	}
	minor := strconv.Itoa(v.Minor)
	if len(minor) == 1 {
		minor = "0" + minor
	}
	update := strconv.Itoa(v.Update)
	if len(update) == 1 {
		update = "00" + update
	} else if len(update) == 1 {
		minor = "0" + update
	}
	return major + "." + minor + "." + update
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
