package models

import (
	"fmt"
)

type Group struct {
	Name         string
	Level        int
	Amount       int
	Competencies []Competency
}

func (g Group) String() string {
	s := g.Name

	if g.Level > 0 {
		s += fmt.Sprintf(": level %d", g.Level)
	}

	if g.Amount == 0 {
		s += " (any of)"
	} else {
		s += fmt.Sprintf(" (%d of)", g.Amount)
	}

	return s
}

func (g Group) Key() string {
	return fmt.Sprintf("%s-%d", g.Name, g.Level)
}
