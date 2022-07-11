package patterns

import "regexp"

var (
	RoleTitle             *regexp.Regexp = regexp.MustCompile(`^(?:# ){1}(.*)`)
	Skills                *regexp.Regexp = regexp.MustCompile(`(?m)^(?:(?P<amount>\d+)(?: of ))?(?P<name>[^:\r\n]+)+:?(?P<level>\d+)?$`)
	InheritNode           *regexp.Regexp = regexp.MustCompile(`<inherit doc="([^"]+)"\/>`)
	SkillsContainer       *regexp.Regexp = regexp.MustCompile(`<skills>[\r\n]([^<]+)<\/skills>`)
	CompetencyHeaderSplit *regexp.Regexp = regexp.MustCompile(`(?m)^# `)
	CompetencyHeader      *regexp.Regexp = regexp.MustCompile(`(?m)^# (?:(?P<group>[^:\r\n]+)?(?:: ))?(?P<name>.*?)+(?: Level )?(?P<level>\d+)?$`)
)
