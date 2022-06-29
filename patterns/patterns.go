package patterns

import "regexp"

var (
	Group           *regexp.Regexp = regexp.MustCompile(`^([0-9]+) of (.+)$`)
	Competency      *regexp.Regexp = regexp.MustCompile(`^([^:\n\r]*)(:)?(\d+)?$`)
	InheritsNode    *regexp.Regexp = regexp.MustCompile(`<inherit doc="([^"]+)"/>`)
	SkillsContainer *regexp.Regexp = regexp.MustCompile(`(?s)<skills>([^<]+)</skills>`)
)
