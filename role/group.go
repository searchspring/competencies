package role

import (
	"fmt"
)

// GroupAmount is an int type whose String method returns "any" if its == 0.
// In this way, if a group doesn't have an amount requirement it will be represented correctly still.
type GroupAmount int

func (a GroupAmount) String() string {
	if a == 0 {
		return "any"
	}
	return fmt.Sprintf("%d", a)
}

type Group struct {
	Name   string
	Amount *GroupAmount
	Level  int
}

func (g Group) Key() string {
	return fmt.Sprintf("%s-%d", g.Name, g.Level)
}
