package main

import (
	"regexp"
)

type Handler struct {
	Pattern string            `yaml:"pattern"`
	Script  string            `yaml:"script"`
	Stream  bool              `yaml:"stream"`
	Env     map[string]string `yaml:"env"`

	re *regexp.Regexp
}

func (h *Handler) Match(input string) (bool, []string) {
	matches := h.re.FindStringSubmatch(input)
	n := len(matches)
	if n > 1 {
		matches = matches[1:]
	}
	return n > 0, matches
}
