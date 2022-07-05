package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/codeallthethingz/competencies/extensions"
	"github.com/codeallthethingz/competencies/models/competency"
	"github.com/codeallthethingz/competencies/models/role"
)

func main() {
	_, err := getCompetencies()
	if err != nil {
		log.Fatal(err)
	}

	_, err = getRoles()
	if err != nil {
		log.Fatal(err)
	}

	// if err := doc.Generate(competencies, roles); err != nil {
	// 	log.Fatal(err)
	// }
}

func getCompetencies() (map[string]*competency.Competency, error) {
	files, err := ioutil.ReadDir(competency.DirPath)
	if err != nil {
		return nil, err
	}

	competencies := map[string]*competency.Competency{}
	for _, file := range files {
		fn := file.Name()
		if !strings.HasSuffix(fn, extensions.Markdown) {
			continue
		}

		newCompetencies, err := competency.New(fn)
		if err != nil {
			return nil, err
		}
		for _, c := range newCompetencies {
			competencies[c.Key()] = &c
		}
	}

	return competencies, nil
}

func getRoles() (map[string]*role.Role, error) {
	files, err := ioutil.ReadDir(role.DirPath)
	if err != nil {
		return nil, err
	}

	roles := map[string]*role.Role{}
	for _, file := range files {
		fn := file.Name()

		// TODO: REMOVE
		if fn != "engineering-developer-3.md" {
			continue
		}

		if !strings.HasSuffix(fn, extensions.Markdown) {
			continue
		}

		role, err := role.New(fn)
		if err != nil {
			return nil, err
		}
		roles[role.Name] = role
	}

	return roles, nil
}
