package competency

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// TODO evaluate what should and shouldn't be private

const (
	dirPath string = "./competencies"
)

// competencies map singleton set by GetCompetencies()
var competencies map[string]*Competency

type Competency struct {
	Name        string
	Description string
	Filename    string
	Group       string
}

func New(filename string) (*Competency, error) {
	contents, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	name, group := parseTitle(contents)
	description := getDescription(contents)

	return &Competency{
		Name:        name,
		Description: description,
		Filename:    filename,
		Group:       group,
	}, nil
}

func GetCompetencies() (map[string]*Competency, error) {
	if competencies != nil {
		return competencies, nil
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	competencies = map[string]*Competency{}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		mdFilename := file.Name()

		competency, err := New(mdFilename)
		if err != nil {
			return nil, err
		}
		competencies[competency.Name] = competency
	}

	return competencies, nil
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dirPath, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getDescription(contents string) string {
	afterTitle := strings.Join(strings.Split(contents, "\n")[1:], "\n")
	beforeNextHeader := strings.Split(afterTitle, "#")[0]
	return strings.TrimSpace(beforeNextHeader)

}

func parseTitle(contents string) (string, string) {
	titleLine := strings.TrimSpace(strings.Split(contents, "\n")[0][1:])
	titleWithoutLevel := strings.Split(titleLine, "Level")[0]

	titleParts := strings.Split(titleWithoutLevel, ":")
	name := ""
	group := ""
	if len(titleParts) == 1 {
		name = strings.ToLower(strings.TrimSpace(titleParts[0]))
		group = ""
	} else {
		name = strings.ToLower(strings.TrimSpace(titleParts[1]))
		group = strings.ToLower(strings.TrimSpace(titleParts[0]))
	}

	name = strings.TrimSuffix(name, "-")
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	return name, group
}
