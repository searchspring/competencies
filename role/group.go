package role

import "fmt"

type Group struct {
	Name   string
	Amount int
	Level  int
}

func (g Group) Key() string {
	return fmt.Sprintf("%s-%d", g.Name, g.Level)
}
