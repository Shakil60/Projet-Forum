package helper

// Chargement et rendu des gabarits HTML cote serveur.

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

// Fonctions utilitaires disponibles dans les templates (images TMDB, texte, etc.).
var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"tmdbPoster": func(path string) string {
		if path == "" {
			return "/static/img/no-poster.svg"
		}
		return "https://image.tmdb.org/t/p/w500" + path
	},
	"tmdbProfile": func(path string) string {
		if path == "" {
			return "/static/img/no-profile.svg"
		}
		return "https://image.tmdb.org/t/p/w300" + path
	},
	"tmdbBackdrop": func(path string) string {
		if path == "" {
			return ""
		}
		return "https://image.tmdb.org/t/p/w780" + path
	},
	"truncate": func(text string, max int) string {
		if len(text) <= max {
			return text
		}
		return text[:max] + "..."
	},
	"mediaTitle": func(title, name string) string {
		if title != "" {
			return title
		}
		return name
	},
	"mediaDate": func(releaseDate, firstAirDate string) string {
		if releaseDate != "" {
			return releaseDate
		}
		return firstAirDate
	},
}

// Pre-charge toutes les pages avec leur layout et les met en cache.
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

// Affiche la page demandee avec ses donnees, ou renvoie une erreur si elle echoue.
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
