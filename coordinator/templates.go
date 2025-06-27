package coordinator

import (
	"html/template"
	"log"
)

var Templates *template.Template

func InitTemplates() {
	Templates = template.Must(template.ParseGlob("static/templates/*.html"))
	log.Println("Templates loaded successfully")
}
