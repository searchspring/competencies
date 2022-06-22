package role

import "fmt"

type Skill struct {
	Name     string
	Level    int
	Filename string
}

func (s Skill) Key() string {
	return fmt.Sprintf("%s-%d", s.Name, s.Level)
}
