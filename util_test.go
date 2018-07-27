package main

import (
	"testing"
)

func TestCleanupText(t *testing.T) {
	examples := map[string]string{
		"<http://foo.com|foo.com>":                                                 "foo.com",
		"sample link <https://google.com|google.com> and <http://bar.com|bar.com>": "sample link google.com and bar.com",
	}

	for input, expected := range examples {
		result := cleanupText(input)
		if result != expected {
			t.Fatalf("Expected %s but got %s", expected, result)
		}
	}
}
