package controllers

// Gere les pages du catalogue (films, series, personnes) via l'API TMDB.

import (
	"forum/helper"
	"forum/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CatalogController struct {
	tmdbService *services.TMDBService
	renderer    *helper.Renderer
}

func InitCatalogController(tmdbService *services.TMDBService, renderer *helper.Renderer) *CatalogController {
	return &CatalogController{
		tmdbService: tmdbService,
		renderer:    renderer,
	}
}

// Affiche une page d'erreur quand la cle TMDB n'est pas configuree.
func (c *CatalogController) renderConfigError(w http.ResponseWriter, r *http.Request) {
	data := baseData(r)
	data["Error"] = c.tmdbService.ConfigError().Error()
	c.renderer.Render(w, http.StatusServiceUnavailable, "catalog_config.html", data)
}

// Affiche l'accueil du catalogue avec les films et series populaires.
func (c *CatalogController) Home(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	page := parsePage(r)
	movies, moviesErr := c.tmdbService.GetPopularMovies(page)
	if moviesErr != nil {
		c.renderError(w, r, moviesErr)
		return
	}

	series, seriesErr := c.tmdbService.GetPopularSeries(page)
	if seriesErr != nil {
		c.renderError(w, r, seriesErr)
		return
	}

	data := baseData(r)
	data["Movies"] = movies.Results
	data["Series"] = series.Results
	data["Page"] = page
	data["HasMore"] = page < movies.TotalPages || page < series.TotalPages
	data["Search"] = strings.TrimSpace(r.URL.Query().Get("q"))
	c.renderer.Render(w, http.StatusOK, "catalog.html", data)
}

// Liste les films populaires ou les resultats de recherche.
func (c *CatalogController) Movies(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	page := parsePage(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var (
		items    interface{}
		hasMore  bool
		total    int
		listErr  error
	)

	if query != "" {
		result, err := c.tmdbService.SearchMovies(query, page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	} else {
		result, err := c.tmdbService.GetPopularMovies(page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	}

	if listErr != nil {
		c.renderError(w, r, listErr)
		return
	}

	data := baseData(r)
	data["Items"] = items
	data["Page"] = page
	data["HasMore"] = hasMore
	data["Search"] = query
	data["Total"] = total
	data["ListType"] = "movies"
	c.renderer.Render(w, http.StatusOK, "catalog_list.html", data)
}

// Liste les series populaires ou les resultats de recherche.
func (c *CatalogController) Series(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	page := parsePage(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var (
		items   interface{}
		hasMore bool
		total   int
		listErr error
	)

	if query != "" {
		result, err := c.tmdbService.SearchSeries(query, page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	} else {
		result, err := c.tmdbService.GetPopularSeries(page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	}

	if listErr != nil {
		c.renderError(w, r, listErr)
		return
	}

	data := baseData(r)
	data["Items"] = items
	data["Page"] = page
	data["HasMore"] = hasMore
	data["Search"] = query
	data["Total"] = total
	data["ListType"] = "series"
	c.renderer.Render(w, http.StatusOK, "catalog_list.html", data)
}

// Liste les personnes populaires ou les resultats de recherche.
func (c *CatalogController) People(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	page := parsePage(r)
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var (
		items   interface{}
		hasMore bool
		total   int
		listErr error
	)

	if query != "" {
		result, err := c.tmdbService.SearchPeople(query, page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	} else {
		result, err := c.tmdbService.GetPopularPeople(page)
		listErr = err
		items = result.Results
		hasMore = page < result.TotalPages
		total = result.TotalResults
	}

	if listErr != nil {
		c.renderError(w, r, listErr)
		return
	}

	data := baseData(r)
	data["Items"] = items
	data["Page"] = page
	data["HasMore"] = hasMore
	data["Search"] = query
	data["Total"] = total
	data["ListType"] = "people"
	c.renderer.Render(w, http.StatusOK, "catalog_list.html", data)
}

// Recherche globale sur les films, series et personnes.
func (c *CatalogController) Search(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
		return
	}

	page := parsePage(r)
	results, err := c.tmdbService.SearchAll(query, page)
	if err != nil {
		c.renderError(w, r, err)
		return
	}

	data := baseData(r)
	data["Results"] = results
	data["Page"] = page
	c.renderer.Render(w, http.StatusOK, "catalog_search.html", data)
}

// Affiche la fiche detaillee d'un film avec son casting.
func (c *CatalogController) ShowMovie(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "Identifiant film invalide", http.StatusBadRequest)
		return
	}

	movie, movieErr := c.tmdbService.GetMovie(id)
	if movieErr != nil {
		c.renderError(w, r, movieErr)
		return
	}

	credits, creditsErr := c.tmdbService.GetMovieCredits(id)
	if creditsErr != nil {
		c.renderError(w, r, creditsErr)
		return
	}

	data := baseData(r)
	data["Movie"] = movie
	data["Cast"] = services.TopCast(credits.Cast, 12)
	data["Directors"] = services.FilterDirectors(credits.Crew)
	c.renderer.Render(w, http.StatusOK, "catalog_movie.html", data)
}

// Affiche la fiche detaillee d'une serie avec son casting.
func (c *CatalogController) ShowSeries(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "Identifiant serie invalide", http.StatusBadRequest)
		return
	}

	show, showErr := c.tmdbService.GetSeries(id)
	if showErr != nil {
		c.renderError(w, r, showErr)
		return
	}

	credits, creditsErr := c.tmdbService.GetSeriesCredits(id)
	if creditsErr != nil {
		c.renderError(w, r, creditsErr)
		return
	}

	data := baseData(r)
	data["Show"] = show
	data["Cast"] = services.TopCast(credits.Cast, 12)
	data["Directors"] = services.FilterDirectors(credits.Crew)
	c.renderer.Render(w, http.StatusOK, "catalog_show.html", data)
}

// Affiche la fiche detaillee d'une personne.
func (c *CatalogController) ShowPerson(w http.ResponseWriter, r *http.Request) {
	if !c.tmdbService.IsConfigured() {
		c.renderConfigError(w, r)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		http.Error(w, "Identifiant personne invalide", http.StatusBadRequest)
		return
	}

	person, personErr := c.tmdbService.GetPerson(id)
	if personErr != nil {
		c.renderError(w, r, personErr)
		return
	}

	data := baseData(r)
	data["Person"] = person
	c.renderer.Render(w, http.StatusOK, "catalog_person.html", data)
}

// Affiche une page d'erreur en cas de probleme avec l'API TMDB.
func (c *CatalogController) renderError(w http.ResponseWriter, r *http.Request, err error) {
	data := baseData(r)
	data["Error"] = err.Error()
	c.renderer.Render(w, http.StatusBadGateway, "catalog_error.html", data)
}
