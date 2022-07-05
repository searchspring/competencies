package role

import (
	"log"
	"strconv"
	"strings"

	"github.com/codeallthethingz/competencies/clients/file"
	"github.com/codeallthethingz/competencies/models/competency"
	"github.com/codeallthethingz/competencies/models/group"
	"github.com/codeallthethingz/competencies/patterns"
)

const DirPath string = "./roles"

type Role struct {
	Name                   string              `json:"name,omitempty"`
	Level                  *int                `json:"level,omitempty"`
	Skills                 Skills              `json:"skills,omitempty"`
	InheritedRoles         []*Role             `json:"inheritedRoles,omitempty"`
	InheritedRoleFilenames map[string]struct{} `json:"inheritedRoleFilenames,omitempty"`
	Markdown               string              `json:"markdown,omitempty"`
	Filename               string              `json:"filename,omitempty"`
}

type Skills struct {
	Competencies []competency.Competency
	Groups       []group.Group
}

func New(filename string) (*Role, error) {
	return new(filename, false)
}

func new(filename string, inherited bool) (*Role, error) {
	role := &Role{
		Skills: Skills{
			Competencies: []competency.Competency{},
			Groups:       []group.Group{},
		},
		InheritedRoles:         []*Role{},
		InheritedRoleFilenames: map[string]struct{}{},
		Filename:               filename,
	}
	if err := role.build(inherited); err != nil {
		return nil, err
	}
	return role, nil
}

func (role *Role) build(inherited bool) error {
	contents, err := file.Read(DirPath, role.Filename)
	if err != nil {
		return err
	}
	role.Markdown = contents
	role.Name = getName(contents)
	role.Skills = getSkills(contents)

	if !inherited {
		if err := role.getInherited(contents); err != nil {
			return err
		}

		// b, err := json.MarshalIndent(role, "", "  ")
		// if err != nil {
		// 	panic(err)
		// }
		// log.Println(string(b))
	}

	// dedupe

	return nil
}

func (role *Role) getInherited(contents string) error {
	match := patterns.InheritNode.FindStringSubmatch(contents)
	if len(match) != 2 || match[1] == "" {
		return nil
	}
	inheritedFilename := match[1]
	log.Println("inheriting from", inheritedFilename)

	if _, ok := role.InheritedRoleFilenames[inheritedFilename]; !ok {
		role.InheritedRoleFilenames[inheritedFilename] = struct{}{}
		inheritedRole, err := new(inheritedFilename, true)
		if err != nil {
			return err
		}
		role.InheritedRoles = append(role.InheritedRoles, inheritedRole)
	}

	inheritedContents, err := file.Read(DirPath, inheritedFilename)
	if err != nil {
		return err
	}

	if err := role.getInherited(inheritedContents); err != nil {
		return err
	}

	return nil
}

func getName(contents string) string {
	return patterns.RoleTitle.FindStringSubmatch(contents)[1]
}

func getSkills(contents string) Skills {
	skills := Skills{Competencies: []competency.Competency{}, Groups: []group.Group{}}

	skillsString := patterns.SkillsContainer.FindStringSubmatch(string(contents))
	if len(skillsString) != 2 {
		return skills
	}

	skillsMatches := patterns.Skills.FindAllStringSubmatch(skillsString[1], -1)
	for _, match := range skillsMatches {
		if len(match) != 4 {
			continue
		} else if strings.TrimSpace(match[0]) == "" {
			continue
		}

		var amount *int
		if match[1] != "" {
			if amnt, err := strconv.Atoi(match[1]); err == nil {
				amount = &amnt
			}
		}
		name := strings.ToLower(strings.TrimSpace(match[2]))
		level, _ := strconv.Atoi(match[3])

		if amount != nil {
			skills.Groups = append(skills.Groups, *group.New(name, level, amount))
		} else {
			skills.Competencies = append(skills.Competencies, competency.Competency{Name: name, Level: level})
		}
	}

	return skills
}
