package doc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codeallthethingz/competencies/extensions"
	"github.com/codeallthethingz/competencies/models/role"
)

const (
	dirPath  string      = "./docs"
	docPerms os.FileMode = 0644
)

func Generate(roles map[string]*role.Role) error {
	for _, role := range roles {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		appData, _ := ioutil.ReadFile("app.js")
		styleData, _ := ioutil.ReadFile("docs/style.css")
		options := ""
		html := `
			<html>
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
		html += string(tailwind([]byte(role.HTML())))
		html += `</div></body><script>let options='` + options + `';` + string(appData) + `</script></html>`

		htmlFilename := strings.Replace(role.Filename, extensions.Markdown, extensions.HTML, 1)
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", dirPath, htmlFilename), []byte(html), docPerms); err != nil {
			return err
		}
	}

	return nil
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

func singleLine(text []byte) string {
	result := string(text)
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}
