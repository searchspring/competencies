package doc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/codeallthethingz/competencies/extensions"
	"github.com/codeallthethingz/competencies/models/competency"
	"github.com/codeallthethingz/competencies/models/group"
	"github.com/codeallthethingz/competencies/models/role"
	"github.com/codeallthethingz/competencies/patterns"
	"github.com/russross/blackfriday/v2"
)

const (
	dirPath       string      = "./docs"
	docPerms      os.FileMode = 0644
	competencyURL string      = "https://github.com/searchspring/competencies/blob/master/competencies/" // TODO: build via config
)

func Generate(competencies map[string]*competency.Competency, roles map[string]*role.Role) error {
	for _, role := range roles {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		html, err := generateHTML(role, competencies)
		if err != nil {
			return err
		}

		htmlFilename := strings.Replace(role.Filename, extensions.Markdown, extensions.HTML, 1)
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", dirPath, htmlFilename), []byte(html), docPerms); err != nil {
			return err
		}
	}

	return nil
}

func generateHTML(role *role.Role, competencies map[string]*competency.Competency) (string, error) {
	markdown := role.Markdown

	skillsMatch := patterns.Skills.FindStringSubmatch(markdown)
	roleHTML, err := generateRoleHTML(role, competencies)
	if err != nil {
		return "", err
	}

	markdown = strings.Replace(markdown, skillsMatch[0], roleHTML, -1)

	html := blackfriday.Run([]byte(markdown))
	html = tailwind(html)

	appData, _ := ioutil.ReadFile("app.js")
	styleData, _ := ioutil.ReadFile("docs/style.css")

	preContent := `<html>
		<head>
		<title>` + role.Name + `</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<style type='text/css'>` + singleLine(styleData) + `</style>
		<link rel="icon" href="seedling.png" type="image/png">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.1/jquery.min.js" integrity="sha256-CSXorXvZcTkaix6Yvo6HppcZGetbYMGWSFlBw8HfCJo=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-cookie/1.4.1/jquery.cookie.min.js" integrity="sha256-1A78rJEdiWTzco6qdn3igTBv9VupN3Q1ozZNTR4WE/Y=" crossorigin="anonymous"></script>
		<script src="https://apis.google.com/js/api.js"></script>
		<link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon"> 
		<link href="https://unpkg.com/tailwindcss@^1.0/dist/tailwind.min.css" rel="stylesheet">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.11.2/css/all.min.css" integrity="sha256-+N4/V/SbAFiW1MPBCXnfnP9QSN3+Keu+NlB+0ev/YKQ=" crossorigin="anonymous" />
		</head>
		<body>
		<div id='content'>
		`

	options := ""
	postContent := `</div></body><script>\nlet options='` + options + `';\n` + string(appData) + `\n</script></html>`

	return preContent + string(html) + postContent, nil
}

func generateRoleHTML(role *role.Role, competencies map[string]*competency.Competency) (string, error) {
	html := `<div class="skill-group p-4 bg-white shadow-xl mb-4 rounded-lg">`
	// for _, skill := range role.GetSkills() {
	// 	html += getSkillHTML(skill, competencies)
	// }
	// for _, group := range role.GetGroups() {
	// 	html += getGroupHTML(group, competencies)
	// }
	// html += "</div>"

	// for _, inherited := range role.GetInheritedRoles() {
	// 	html += `<div class="skill-group p-4 bg-white shadow-xl mb-4 rounded-lg">`
	// 	html += `<h4 class="px-2 text-l mt-2">` + inherited.Title + `</h4>`
	// 	for _, skill := range inherited.GetSkills() {
	// 		html += getSkillHTML(skill, competencies)
	// 	}
	// 	for _, group := range inherited.GetGroups() {
	// 		html += getGroupHTML(group, competencies)
	// 	}
	// 	html += "</div>"
	// }

	return html, nil
}

func getSkillHTML(skill competency.Competency, competencies map[string]*competency.Competency) string {
	html := ""
	classes := ""
	href := createHREF(skill.Name)
	drive := " <a href=\"javascript:;\" title=\"add this competency to the google sheet for tracking\" style=\"display:none\" class=\"drive-link hover:opacity-75\"><i class=\"fas hover:opacity-75 ml-1 fa-plus\"></i></a>"
	if _, ok := competencies[skill.Name]; !ok {
		classes += "missing"
		href = "https://github.com/searchspring/competencies/new/master/competencies"
	}
	github := "<a href=\"" + href + "\" title=\"go to competency github page\" class=\"github-link\" target=\"_blank\"><i class=\"fab hover:opacity-75 fa-github\"></i></a> "
	classes += " " + name2Id(skill.Name)
	classes += " competency inline-block rounded-full bg-gray-300 p-1 px-2 mr-2 mb-2 text-xs whitespace-no-wrap"

	html += "<span id=\"" + name2Id(skill.Name) + "\" level=\"" + strconv.Itoa(skill.Level) + "\" class=\"" + classes + "\">" + github + nameLevelAndGroup(skill.Name, skill.Level, nil) + drive + "</span>"
	return html
}

func getGroupHTML(group group.Group, competencies map[string]*competency.Competency) string {
	// html := `<table class="group mt-4"><tr><td valign="top"><span class="group-heading text-sm pr-2 whitespace-no-wrap">` + nameLevelAndGroup(group.Name, group.Level, group.Amount) + `</span></td><td class="group" valign="top"> `
	// for n, c := range competencies {
	// 	if c.Group == group.Name {
	// 		html += getSkillHTML(role.Skill{Name: n, Level: group.Level, Filename: c.Filename}, competencies)
	// 	}
	// }
	// html += `</td></tr></table>`
	// return html
	return ""
}

func tailwind(html []byte) []byte {
	htmlString := string(html)
	htmlString = strings.ReplaceAll(htmlString, "<h1>", `<h1 class="whitespace-no-wrap top-0 left-0 fixed w-full block opacity-90 bg-white p-2 px-8 border-b-2 text-lg mb-4"><img class="w-6 inline-block mr-3" src="seedling.png">`)
	htmlString = strings.ReplaceAll(htmlString, "<h2>", `<h2 class="px-2 text-2xl mt-4">`)
	htmlString = strings.ReplaceAll(htmlString, "<h3>", `<h3 class="px-2 text-xl mt-2">`)
	htmlString = strings.ReplaceAll(htmlString, "<h4>", `<h4 class="px-2 text-l mt-2">`)
	htmlString = strings.ReplaceAll(htmlString, "<p>", `<p style="width:50rem" class="px-2">`)
	return []byte(htmlString)
}

func name2Id(name string) string {
	return "c-" + strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", ""), "-", ""))
}

func nameLevelAndGroup(name string, level int, groupAmount *int) string {
	lvlstr := ""
	if level > 1 {
		lvlstr = fmt.Sprintf(": level %d", level)
	}

	grpstr := ""
	if groupAmount != nil {
		grpstr = fmt.Sprintf(" (%d of)", groupAmount)
	}

	return fmt.Sprintf("%s%s%s", strings.ToLower(strings.TrimSpace(name)), lvlstr, grpstr)
}

func singleLine(text []byte) string {
	result := string(text)
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}

func createHREF(name string) string {
	return competencyURL + cleanFile(name)
}

func cleanFile(file string) string {
	file = strings.ReplaceAll(file, " - ", "-")
	file = strings.ReplaceAll(strings.TrimSpace(file), " ", "-")
	if strings.HasSuffix(file, ":2") || strings.HasSuffix(file, ":3") || strings.HasSuffix(file, ":4") {
		file = file[0 : len(file)-2]
	}
	return strings.ToLower(file) + ".md"
}
