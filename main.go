package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/codeallthethingz/competencies/role"
)

func main() {
	files, err := ioutil.ReadDir(role.RolesDirPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(role.DocsDirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		// TODO: REMOVE
		if file.Name() != "engineering-developer-3.md" {
			continue
		}

		if !strings.HasSuffix(file.Name(), role.MarkdownExtension) {
			continue
		}
		mdFilename := file.Name()
		log.Println(mdFilename)

		role, err := role.New(mdFilename)
		if err != nil {
			log.Fatal(err)
		}

		for filename, role := range role.Inherited {
			log.Println(filename, role.Title, role.Skills, role.Groups)
		}
	}
}
