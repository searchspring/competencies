package competency

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeallthethingz/competencies/clients/file"
	"github.com/codeallthethingz/competencies/patterns"
)

const (
	DirPath  string = "./competencies"
	mdHeader        = "# "
)

type Competency struct {
	Filename  string `json:"filename,omitempty"`
	Markdown  string `json:"-"`
	GroupName string `json:"group,omitempty"`
	Name      string `json:"name,omitempty"`
	Level     int    `json:"level,omitempty"`
	Missing   bool   `json:"Missing"`
}

// TODO: this does work that may not belong in new and it should only return one
func New(filename string) ([]Competency, error) {
	allContents, err := file.Read(DirPath, filename)
	if err != nil {
		return nil, err
	}

	competencies := []Competency{}
	for _, splitContents := range patterns.CompetencyHeaderSplit.Split(allContents, -1) {
		trimmedSplitContents := strings.TrimSpace(splitContents)
		if trimmedSplitContents == "" {
			continue
		}
		contents := mdHeader + trimmedSplitContents

		match := patterns.CompetencyHeader.FindStringSubmatch(contents)
		if len(match) != 4 {
			continue
		} else if strings.TrimSpace(match[0]) == "" {
			continue
		}

		level, _ := strconv.Atoi(strings.TrimSpace(match[3]))

		competency := Competency{
			Filename:  filename,
			Markdown:  contents,
			GroupName: strings.ToLower(strings.TrimSpace(match[1])),
			Name:      strings.ToLower(strings.TrimSpace(match[2])),
			Level:     level,
		}

		competencies = append(competencies, competency)
	}

	return competencies, nil
}

func (c Competency) Key() string {
	return fmt.Sprintf("%s-%d", c.Name, c.Level)
}

func (c Competency) HTML() string {
	html := ""
	classes := ""
	href := "https://github.com/searchspring/competencies/blob/master/competencies/" + c.Filename
	drive := " <a href=\"javascript:;\" title=\"add this competency to the google sheet for tracking\" style=\"display:none\" class=\"drive-link hover:opacity-75\"><i class=\"fas hover:opacity-75 ml-1 fa-plus\"></i></a>"
	if c.Missing {
		classes += "missing"
		href = "https://github.com/searchspring/competencies/new/master/competencies"
	}
	github := "<a href=\"" + href + "\" title=\"go to competency github page\" class=\"github-link\" target=\"_blank\"><i class=\"fab hover:opacity-75 fa-github\"></i></a> "
	id := "c-" + strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(c.Name, " ", ""), "-", ""))
	classes += " " + id
	classes += " competency inline-block rounded-full bg-gray-300 p-1 px-2 mr-2 mb-2 text-xs whitespace-no-wrap"
	html += "<span id=\"" + id + "\" level=\"" + strconv.Itoa(c.Level) + "\" class=\"" + classes + "\">" + github + c.String() + drive + "</span>"
	return html
}

func (c Competency) String() string {
	s := c.Name
	if c.Level > 0 {
		s += fmt.Sprintf(": level %d", c.Level)
	}
	return s
}

func EqualLevel(levelA, levelB int) bool {
	return levelA == levelB || levelA == 1 && levelB == 0 || levelA == 0 && levelB == 1
}
