package group

import (
	"fmt"

	"github.com/codeallthethingz/competencies/models/competency"
)

type Group struct {
	Name         string
	Level        int
	Amount       *int
	Competencies []competency.Competency
}

func New(name string, level int, amount *int) *Group {
	return &Group{
		Name:         name,
		Level:        level,
		Amount:       amount,
		Competencies: []competency.Competency{},
	}
	// TODO: Get competencies from markdown files
}

// func (g Group) String() string {
// 	s := g.Name

// 	if g.Level > 0 {
// 		s += fmt.Sprintf(": level %d", g.Level)
// 	}

// 	if g.Amount == 0 {
// 		s += " (any of)"
// 	} else {
// 		s += fmt.Sprintf(" (%d of)", g.Amount)
// 	}

// 	return s
// }

func (g Group) Key() string {
	return fmt.Sprintf("%s-%d", g.Name, g.Level)
}
