package patterns

import "regexp"

var (
	Group    *regexp.Regexp = regexp.MustCompile(`^([0-9]+) of (.+)$`)
	Inherits *regexp.Regexp = regexp.MustCompile(`<inherit doc="([^"]+)"/>`)
	Skills   *regexp.Regexp = regexp.MustCompile(`(?s)<skills>([^<]+)</skills>`)
)
