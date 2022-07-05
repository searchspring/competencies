package patterns

import "regexp"

var (
	RoleTitle                    *regexp.Regexp = regexp.MustCompile(`^(?:# ){1}(.*)`)
	Skills                       *regexp.Regexp = regexp.MustCompile(`(?m)^(\d+)?(?: of )?([^:\r\n]+)+(?::)?(\d+)?$`)
	InheritNode                  *regexp.Regexp = regexp.MustCompile(`<inherit doc="([^"]+)"/>`)
	SkillsContainer              *regexp.Regexp = regexp.MustCompile(`<skills>[\r\n]([^<]+)</skills>`)
	CompetenciesContentsSplitter *regexp.Regexp = regexp.MustCompile(`(?m)^# `)
)
