package role

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeallthethingz/competencies/clients/file"
	"github.com/codeallthethingz/competencies/models/competency"
	"github.com/codeallthethingz/competencies/models/group"
	"github.com/codeallthethingz/competencies/patterns"
	"github.com/russross/blackfriday/v2"
	log "github.com/sirupsen/logrus"
)

const DirPath string = "./roles"

var missingCompetenciesMessages map[string]struct{} = map[string]struct{}{}

type Role struct {
	Name                   string                  `json:"name,omitempty"`
	Level                  *int                    `json:"level,omitempty"`
	Skills                 Skills                  `json:"skills,omitempty"`
	InheritedRoles         []*Role                 `json:"inheritedRoles,omitempty"`
	InheritedRoleFilenames map[string]struct{}     `json:"inheritedRoleFilenames,omitempty"`
	Markdown               string                  `json:"-"`
	Filename               string                  `json:"filename,omitempty"`
	competencies           []competency.Competency // master competencies slice
}

type Skills struct {
	Competencies []competency.Competency
	Groups       []group.Group
}

func (skills Skills) HTML() string {
	html := `<div class="skill-group p-4 bg-white shadow-xl mb-4 rounded-lg">`
	for _, competency := range skills.Competencies {
		html += competency.HTML()
	}
	for _, group := range skills.Groups {
		html += group.HTML()
	}
	html += `</div>`
	return html
}

func New(filename string, competencies []competency.Competency) (*Role, error) {
	inherited := false
	return new(filename, competencies, inherited)
}

func new(filename string, competencies []competency.Competency, inherited bool) (*Role, error) {
	role := &Role{
		Skills: Skills{
			Competencies: []competency.Competency{},
			Groups:       []group.Group{},
		},
		InheritedRoles:         []*Role{},
		InheritedRoleFilenames: map[string]struct{}{},
		Filename:               filename,
		competencies:           competencies,
	}

	contents, err := file.Read(DirPath, role.Filename)
	if err != nil {
		return nil, err
	}
	role.Markdown = contents
	role.getName(contents)
	role.getSkills(contents)

	if !inherited {
		if err := role.getInherited(contents); err != nil {
			return nil, err
		}
	}

	// TODO: dedupe inherited skills, roll up to parent role

	return role, nil
}

func (role *Role) getInherited(contents string) error {
	match := patterns.InheritNode.FindStringSubmatch(contents)
	if len(match) != 2 || match[1] == "" {
		return nil
	}
	inheritedFilename := match[1]
	log.Debugf("inheriting from %s", inheritedFilename)

	if _, ok := role.InheritedRoleFilenames[inheritedFilename]; !ok {
		role.InheritedRoleFilenames[inheritedFilename] = struct{}{}
		inheritedRole, err := new(inheritedFilename, role.competencies, true)
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

func (role *Role) getName(contents string) {
	role.Name = patterns.RoleTitle.FindStringSubmatch(contents)[1]
}

func (role *Role) getSkills(contents string) {
	skills := Skills{Competencies: []competency.Competency{}, Groups: []group.Group{}}

	skillsString := patterns.SkillsContainer.FindStringSubmatch(string(contents))
	if len(skillsString) != 2 {
		return
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
			skills.Groups = append(skills.Groups, *group.New(name, level, amount, role.competencies))
		} else {
			skills.Competencies = append(skills.Competencies, role.getCompetency(name, level))
		}
	}

	role.Skills = skills
}

func (role Role) getCompetency(name string, level int) competency.Competency {
	for _, c := range role.competencies {
		if c.Name == name && competency.EqualLevel(c.Level, level) {
			return c
		}
	}
	tempCompetency := competency.Competency{Name: name, Level: level, Missing: true}
	msg := fmt.Sprintf("missing flat competency: %s", tempCompetency.String())
	if _, ok := missingCompetenciesMessages[msg]; !ok {
		log.Warn(msg)
		missingCompetenciesMessages[msg] = struct{}{}
	}
	return tempCompetency
}

func (role Role) HTML() string {
	md := patterns.SkillsContainer.ReplaceAll([]byte(role.Markdown), []byte(role.Skills.HTML()))

	inherited := ""
	for _, inheritedRole := range role.InheritedRoles {
		inherited += "<h4>" + inheritedRole.Name + "</h4>" + inheritedRole.Skills.HTML()
	}

	md = patterns.InheritNode.ReplaceAll(md, []byte(inherited))
	return string(blackfriday.Run([]byte(md)))
}
