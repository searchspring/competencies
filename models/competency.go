package models

import "fmt"

type Competency struct {
	Name  string
	Level int
}

func (c Competency) Key() string {
	return fmt.Sprintf("%s-%d", c.Name, c.Level)
}
