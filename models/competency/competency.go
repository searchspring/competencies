package competency

import (
	"fmt"
	"strings"

	"github.com/codeallthethingz/competencies/clients/file"
	"github.com/codeallthethingz/competencies/patterns"
)

const DirPath string = "./competencies"

type Competency struct {
	Filename  string
	Markdown  string
	GroupName string
	Name      string
	Level     int
}

func New(filename string) ([]Competency, error) {
	allContents, err := file.Read(DirPath, filename)
	if err != nil {
		return nil, err
	}

	for _, splitContents := range patterns.CompetenciesContentsSplitter.Split(allContents, -1) {
		trimmedSplitContents := strings.TrimSpace(splitContents)
		if trimmedSplitContents == "" {
			continue
		}
		contents := "# " + trimmedSplitContents

		competency := &Competency{Filename: filename, Markdown: contents}
		competency.Name = "asdas" // TODO: CONTINUE FROM HERE
	}

	return nil, nil
}

func (c Competency) Key() string {
	return fmt.Sprintf("%s-%d", c.Name, c.Level)
}
