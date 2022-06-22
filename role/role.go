package role

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	MarkdownExtension string = ".md"
	HtmlExtension     string = ".html"
	RolesDirPath      string = "./roles"
	DocsDirPath       string = "./docs"
)

var (
	groupPattern    *regexp.Regexp = regexp.MustCompile(`^([0-9]+) of (.+)$`)
	inheritsPattern *regexp.Regexp = regexp.MustCompile(`<inherit doc="([^"]+)"/>`)
	skillsPattern   *regexp.Regexp = regexp.MustCompile(`(?s)<skills>([^<]+)</skills>`)
)

type Role struct {
	Title     string
	Filename  string
	Skills    map[string]Skill
	Groups    map[string]Group
	Inherited map[string]*Role
}

func New(filename string) (*Role, error) {
	role := &Role{
		Filename: filename,
		Skills:   map[string]Skill{},
		Groups:   map[string]Group{},
	}
	if err := role.build(); err != nil {
		return nil, err
	}
	return role, nil
}

func (role Role) getSkills() []Skill {
	skills := []Skill{}
	for _, skill := range role.Skills {
		skills = append(skills, skill)
	}
	sort.SliceStable(skills, func(i, j int) bool {
		if skills[i].Name == skills[j].Name {
			return skills[i].Level < skills[j].Level
		}
		return skills[i].Name < skills[j].Name
	})
	return skills
}

func (role Role) getGroups() []Group {
	groups := []Group{}
	for _, group := range role.Groups {
		groups = append(groups, group)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].Name == groups[j].Name {
			return groups[i].Level < groups[j].Level
		}
		return groups[i].Name < groups[j].Name
	})
	return groups
}

func (role *Role) build() error {
	contents, err := readRoleFile(role.Filename)
	if err != nil {
		return err
	}

	role.Title = getTitle(contents)

	skillStrings := readSkillsList(contents)
	role.addSkills(skillStrings)
	role.addSkillGroups(skillStrings)
	if err := role.addInherited(role.Filename); err != nil {
		return err
	}

	return nil
}

func (role *Role) addSkills(skillStrings []string) {
	for _, skillString := range skillStrings {
		if groupPattern.MatchString(skillString) {
			continue
		}

		split := strings.Split(skillString, ":")
		skill := Skill{Name: strings.TrimSpace(split[0]), Level: 0}

		if len(split) > 1 {
			level, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
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
		groupMatches := groupPattern.FindStringSubmatch(skillString)
		if len(groupMatches) <= 1 {
			continue
		}

		split := strings.Split(groupMatches[2], ":")
		group := Group{Name: strings.TrimSpace(split[0])}

		amount, err := strconv.Atoi(groupMatches[1])
		if err != nil {
			panic(err)
		}
		group.Amount = amount

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
	contents, err := readRoleFile(filename)
	if err != nil {
		return err
	}

	matches := inheritsPattern.FindAllStringSubmatch(contents, -1)
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		inheritedFilename := match[1]
		log.Println(inheritedFilename)
		// TODO: this but avoid infinite recursion
	}

	return nil
}

func readRoleFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", RolesDirPath, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getTitle(contents string) string {
	return strings.TrimSpace(strings.Split(contents, "\n")[0][1:])
}

func readSkillsList(contents string) []string {
	skills := []string{}
	match := skillsPattern.FindStringSubmatch(string(contents))
	for _, rawSkillStr := range strings.Split(match[1], "\n") {
		cleanSkillStr := strings.TrimSpace(rawSkillStr)
		if cleanSkillStr == "" {
			continue
		}
		skills = append(skills, cleanSkillStr)
	}
	return skills
}
