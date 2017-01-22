package util

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

// Page

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("pages/"+filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("pages/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func CreatePages(titles ...string) {
	var pg *Page
	for _, title := range titles {
		pg = &Page{Title: title, Body: nil}
		pg.Save()
	}
}

// Templates

type Templates map[string]*template.Template

var Tmpl = Templates{}

func RegisterTemplates(ts ...string) {
	for _, t := range ts {
		Tmpl[t] = template.Must(template.ParseFiles("templates/"+t, "templates/base.html"))
	}
}

func RenderTemplate(w http.ResponseWriter, t string, pg *Page) {
	Tmpl[t].ExecuteTemplate(w, "base", &pg)
}

// Handler

type Handler func(w http.ResponseWriter, req *http.Request)

func TemplateHandler(filename string) Handler {
	return func(w http.ResponseWriter, req *http.Request) {
		pg, _ := LoadPage(string(req.URL.Path[1:]))
		RenderTemplate(w, filename, pg)
	}
}
