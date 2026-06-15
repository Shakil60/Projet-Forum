package helper

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates map[string]*template.Template
}

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
}

func InitRenderer(dir string) *Renderer {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "pages", "*.html"))
	if err != nil {
		log.Fatalf("Erreur chargement des vues - %s", err.Error())
	}

	layout := filepath.Join(dir, "layouts", "base.html")

	for _, page := range pages {
		name := filepath.Base(page)
		ts := template.New(name).Funcs(funcMap)
		ts = template.Must(ts.ParseFiles(layout))
		ts = template.Must(ts.ParseFiles(page))
		cache[name] = ts
	}

	log.Printf("Vues - %d pages chargees", len(cache))
	return &Renderer{templates: cache}
}

func (rd *Renderer) Render(w http.ResponseWriter, status int, page string, data map[string]any) {
	ts, ok := rd.templates[page]
	if !ok {
		http.Error(w, "Vue introuvable : "+page, http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buffer, "base", data); err != nil {
		log.Printf("Erreur rendu vue %s - %s", page, err.Error())
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buffer.WriteTo(w)
}
