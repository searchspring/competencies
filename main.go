package main

import (
	"io/ioutil"
	"strings"

	"github.com/codeallthethingz/competencies/doc"
	"github.com/codeallthethingz/competencies/extensions"
	"github.com/codeallthethingz/competencies/models/competency"
	"github.com/codeallthethingz/competencies/models/role"
	log "github.com/sirupsen/logrus"
)

func main() {
	competencies, err := getCompetencies()
	if err != nil {
		log.Fatal(err)
	}

	roles, err := getRoles(competencies)
	if err != nil {
		log.Fatal(err)
	}

	if err := doc.Generate(roles, competencies); err != nil {
		log.Fatal(err)
	}
}

func getCompetencies() ([]competency.Competency, error) {
	files, err := ioutil.ReadDir(competency.DirPath)
	if err != nil {
		return nil, err
	}

	competencies := []competency.Competency{}
	for _, file := range files {
		fn := file.Name()
		if !strings.HasSuffix(fn, extensions.Markdown) {
			continue
		}

		newCompetencies, err := competency.New(fn)
		if err != nil {
			return nil, err
		}
		competencies = append(competencies, newCompetencies...)
	}

	return competencies, nil
}

func getRoles(competencies []competency.Competency) (map[string]*role.Role, error) {
	files, err := ioutil.ReadDir(role.DirPath)
	if err != nil {
		return nil, err
	}

	roles := map[string]*role.Role{}
	for _, file := range files {
		fn := file.Name()

		if !strings.HasSuffix(fn, extensions.Markdown) {
			continue
		}

		role, err := role.New(fn, competencies)
		if err != nil {
			return nil, err
		}
		roles[role.Name] = role
	}

	return roles, nil
}
