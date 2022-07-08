package group

import (
	"fmt"
	"strings"

	"github.com/codeallthethingz/competencies/models/competency"
	log "github.com/sirupsen/logrus"
)

var missingCompetenciesMessages map[string]struct{} = map[string]struct{}{}

type Group struct {
	Name         string                  `json:"name,omitempty"`
	Level        int                     `json:"level,omitempty"`
	Amount       *int                    `json:"amount,omitempty"`
	Competencies []competency.Competency `json:"competencies,omitempty"`
	competencies []competency.Competency // master competencies slice
}

func New(name string, level int, amount *int, competencies []competency.Competency) *Group {
	g := &Group{
		Name:         strings.ToLower(strings.TrimSpace(name)),
		Level:        level,
		Amount:       amount,
		competencies: competencies,
	}
	g.getCompetencies()
	return g
}

func (g *Group) getCompetencies() {
	g.Competencies = []competency.Competency{}
	added := map[string]struct{}{}
	// Loop through one time adding competencies that exist.
	for _, c := range g.competencies {
		if c.GroupName == g.Name && competency.EqualLevel(g.Level, c.Level) {
			g.Competencies = append(g.Competencies, c)
			added[c.Name] = struct{}{}
		}
	}
	// Loop through a second time adding competencies that are missing.
	// This is separate to avoid marking something as missing early since
	// map iteration is orderless.
	for _, c := range g.competencies {
		if c.GroupName != g.Name {
			continue
		} else if _, ok := added[c.Name]; !ok {
			msg := fmt.Sprintf("missing group competency: %s level %d %s", g.Name, g.Level, c.String())
			if _, ok := missingCompetenciesMessages[msg]; !ok {
				log.Warn(msg)
				missingCompetenciesMessages[msg] = struct{}{}
			}
			g.Competencies = append(g.Competencies, competency.Competency{Name: c.Name, Level: g.Level, GroupName: g.Name, Missing: true})
			added[c.Name] = struct{}{}
		}
	}
}

func (g Group) Key() string {
	return fmt.Sprintf("%s-%d", g.Name, g.Level)
}

func (g Group) HTML() string {
	html := `<table class="group mt-4"><tr><td valign="top"><span class="group-heading text-sm pr-2 whitespace-no-wrap">` + g.String() + `</span></td><td class="group" valign="top"> `
	for _, c := range g.Competencies {
		html += c.HTML()
	}
	html += `</td></tr></table>`
	return html
}

func (g Group) String() string {
	s := g.Name

	if g.Level > 0 {
		s += fmt.Sprintf(": level %d", g.Level)
	}

	if g.Amount != nil {
		if *g.Amount == 0 {
			s += " (any of)"
		} else {
			s += fmt.Sprintf(" (%d of)", *g.Amount)
		}
	}

	return s
}
