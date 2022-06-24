package main

import (
	"log"

	"github.com/codeallthethingz/competencies/competency"
	"github.com/codeallthethingz/competencies/doc"
	"github.com/codeallthethingz/competencies/role"
)

func main() {
	competencies, err := competency.GetCompetencies()
	if err != nil {
		log.Fatal(err)
	}

	roles, err := role.GetRoles()
	if err != nil {
		log.Fatal(err)
	}

	if err := doc.Generate(competencies, roles); err != nil {
		log.Fatal(err)
	}
}
