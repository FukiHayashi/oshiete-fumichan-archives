package render

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
)

func CreateRenderTemplates() multitemplate.Render {
	r := multitemplate.New()

	template_funcs := map[string]interface{}{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	}

	htmls, err := filepath.Glob("./views/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, html := range htmls {
		slice := strings.Split(html, "/")
		r.AddFromFilesFuncs(slice[len(slice)-1], template_funcs, "./views/layout.tmpl", "./views/_header.tmpl",
			"./views/_footer.tmpl", "./views/_messages.tmpl", html)
	}
	return r
}
