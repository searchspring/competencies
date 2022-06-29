package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/codeallthethingz/competencies/extensions"
	"github.com/codeallthethingz/competencies/patterns"
)

// TODO evaluate what should and shouldn't be private

const dirPath string = "./roles"

// roles map singleton set by GetRoles()
var roles map[string]*Role

type Role struct {
	Name      string
	Level     int
	Skills    Skills
	Inherited []*Role
	Markdown  string
	Filename  string
}

type Skills struct {
	Competencies []Competency
	Groups       []Group
}

func New(filename string) (*Role, error) {
	return new(filename, false)
}

func new(filename string, inherited bool) (*Role, error) {
	role := &Role{
		Skills: Skills{
			Competencies: []Competency{},
			Groups:       []Group{},
		},
		Inherited: []*Role{},
		Filename:  filename,
	}
	if err := role.build(inherited); err != nil {
		return nil, err
	}
	return role, nil
}

func GetRoles() (map[string]*Role, error) {
	if roles != nil {
		return roles, nil
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	roles = map[string]*Role{}
	for _, file := range files {
		// TODO: REMOVE
		if file.Name() != "engineering-developer-3.md" {
			continue
		}

		if !strings.HasSuffix(file.Name(), extensions.Markdown) {
			continue
		}
		mdFilename := file.Name()

		role, err := New(mdFilename)
		if err != nil {
			return nil, err
		}
		roles[role.Title] = role
	}

	return roles, nil
}

// func (role Role) GetSkills() []Skill {
// 	skills := []Skill{}
// 	for _, skill := range role.Skills {
// 		skills = append(skills, skill)
// 	}
// 	sort.SliceStable(skills, func(i, j int) bool {
// 		if skills[i].Name == skills[j].Name {
// 			return skills[i].Level < skills[j].Level
// 		}
// 		return skills[i].Name < skills[j].Name
// 	})
// 	return skills
// }

// func (role Role) GetGroups() []Group {
// 	groups := []Group{}
// 	for _, group := range role.Groups {
// 		groups = append(groups, group)
// 	}
// 	sort.SliceStable(groups, func(i, j int) bool {
// 		if groups[i].Name == groups[j].Name {
// 			return groups[i].Level < groups[j].Level
// 		}
// 		return groups[i].Name < groups[j].Name
// 	})
// 	return groups
// }

// func (role Role) GetInheritedRoles() []*Role {
// 	roles := []*Role{}
// 	for _, inheritedRole := range role.Inherited {
// 		roles = append(roles, inheritedRole)
// 	}
// 	sort.SliceStable(roles, func(i, j int) bool {
// 		lenI := len(roles[i].GetSkills()) + len(roles[i].GetGroups())
// 		lenJ := len(roles[j].GetSkills()) + len(roles[j].GetGroups())
// 		return lenI > lenJ
// 	})
// 	return roles
// }

func (role *Role) build(inherited bool) error {
	contents, err := readFile(role.Filename)
	if err != nil {
		return err
	}

	role.Name = GetTitle(contents)
	role.Markdown = contents
	skillStrings := readSkillsList(contents)

	role.parseCompetencies(skillStrings)
	role.addSkillGroups(skillStrings)
	if !inherited {
		// recursively build inherited roles
		if err := role.addInherited(role.Filename); err != nil {
			return err
		}
		role.dedupeSkillsAndGroups()
	}

	return nil
}

func (role *Role) parseCompetencies(skillStrings []string) {
	for _, skillString := range skillStrings {
		if patterns.Group.MatchString(skillString) {
			continue
		}

		split := strings.Split(skillString, ":")
		skill := Skill{Name: strings.ToLower(strings.TrimSpace(split[0])), Level: 0}

		if len(split) > 1 {
			level, err := strconv.Atoi(split[1])
			if err != nil {
				log.Println(err)
				continue
			}
			skill.Level = level
		}

		if _, ok := role.Skills[skill.Key()]; !ok {
			role.Skills[skill.Key()] = skill
		}
	}
}

func (role *Role) addSkillGroups(skillStrings []string) {
	for _, skillString := range skillStrings {
		groupMatches := patterns.Group.FindStringSubmatch(skillString)
		if len(groupMatches) <= 1 {
			continue
		}

		split := strings.Split(groupMatches[2], ":")
		group := Group{Name: strings.ToLower(strings.TrimSpace(split[0]))}

		amount, err := strconv.Atoi(groupMatches[1])
		if err != nil {
			panic(err)
		}
		ga := GroupAmount(amount)
		group.Amount = &ga

		if len(split) > 1 {
			level, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			group.Level = level
		}

		if _, ok := role.Groups[group.Key()]; !ok {
			role.Groups[group.Key()] = group
		}
	}
}

func (role *Role) addInherited(filename string) error {
	log.Println("inheriting from", filename)
	contents, err := readFile(filename)
	if err != nil {
		return err
	}

	matches := patterns.Inherits.FindAllStringSubmatch(contents, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		inheritedFilename := match[1]
		if _, ok := role.Inherited[inheritedFilename]; !ok {
			// create a new inherited role and add it to the main parent role's inherited map
			// calling new with inherited=true will skip calling addInherited on the recursed
			// build call.
			inheritedRole, err := new(inheritedFilename, true)
			if err != nil {
				return err
			}
			role.Inherited[inheritedFilename] = inheritedRole
			// recursively inherit into main parent role
			if err := role.addInherited(inheritedFilename); err != nil {
				return err
			}
		}
	}

	return nil
}

func (role *Role) dedupeSkillsAndGroups() {
	for _, inheritedRole := range role.Inherited {
		for key := range inheritedRole.Skills {
			if _, ok := role.Skills[key]; ok {
				delete(inheritedRole.Skills, key)
			}
		}
		for key := range inheritedRole.Groups {
			if _, ok := role.Groups[key]; ok {
				delete(inheritedRole.Groups, key)
			}
		}
	}
}

func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dirPath, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetTitle(contents string) string {
	return strings.TrimSpace(strings.Split(contents, "\n")[0][1:])
}

func readSkillsList(contents string) []string {
	skills := []string{}
	match := patterns.Skills.FindStringSubmatch(string(contents))
	for _, rawSkillStr := range strings.Split(match[1], "\n") {
		cleanSkillStr := strings.TrimSpace(rawSkillStr)
		if cleanSkillStr == "" {
			continue
		}
		skills = append(skills, cleanSkillStr)
	}
	return skills
}
