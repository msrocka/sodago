package main

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	check := func(s string, major int, minor int, update int) {
		v := ParseVersion(s)
		if v.Major != major || v.Minor != minor || v.Update != update {
			t.Error("failed to parse version", s, "; got", v)
		}
	}
	check("", 0, 0, 0)
	check("01.01.001", 1, 1, 1)
	check("99.88.777", 99, 88, 777)
}

func TestVersionString(t *testing.T) {
	check := func(raw string, expected string) {
		s := ParseVersion(raw).String()
		if s != expected {
			t.Error("failed to stringify version; expected", expected,
				"but got", s)
		}
	}
	check("", "00.00.000")
	check("1.1.1", "01.01.001")
	check("99.88.777", "99.88.777")
}

func TestVersionCompare(t *testing.T) {
	v1 := ParseVersion("1.2.3")
	v2 := ParseVersion("3.2.1")
	if diff := v1.Compare(v2); diff >= 0 {
		t.Error("failed to compare", v1, "with", v2)
	}
	if v2.Compare(v1) <= 0 {
		t.Error("failed to compare", v2, "with", v1)
	}
	if v2.Compare(v2) != 0 {
		t.Error("failed to compare", v2, "with", v2)
	}
}
