package controllers

// Gere les pages du catalogue (films, series, personnes) via l'API TMDB.

import (
	"forum/helper"
	"forum/middleware"
	"forum/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CatalogController struct {
	tmdbService *services.TMDBService
	filmService *services.FilmService
	renderer    *helper.Renderer
}

func InitCatalogController(tmdbService *services.TMDBService, filmService *services.FilmService, renderer *helper.Renderer) *CatalogController {
	return &CatalogController{
		tmdbService: tmdbService,
		filmService: filmService,
		renderer:    renderer,
	}
}

// Affiche une page d'erreur quand la cle TMDB n'est pas configuree.
func (c *CatalogController) renderConfigError(w http.ResponseWriter, r *http.Request) {
	data := baseData(r)
	data["Error"] = c.tmdbService.ConfigError().Error()
	c.renderer.Render(w, http.StatusServiceUnavailable, "catalog_config.html", data)
}

// Ajoute aux donnees les favoris et commentaires d'un film/serie.
func (c *CatalogController) attachFilmSocial(r *http.Request, data map[string]any, tmdbId int, mediaType string) {
	data["TmdbId"] = tmdbId
	data["MediaType"] = mediaType
	data["Comments"] = c.filmService.Comments(tmdbId, mediaType)
	data["FavCount"] = c.filmService.CountFavoris(tmdbId, mediaType)
	if user := middleware.GetUser(r); user != nil {
		data["IsFavori"] = c.filmService.IsFavori(user.Id, tmdbId, mediaType)
	}
}

// Determine le type de media et l'identifiant depuis l'URL.
func mediaParams(r *http.Request) (int, string, bool) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id <= 0 {
		return 0, "", false
	}
	mediaType := "movie"
	if strings.Contains(r.URL.Path, "/series/") {
		mediaType = "tv"
	}
	return id, mediaType, true
}

// Construit l'URL de retour vers la fiche du film/serie.
func filmPath(id int, mediaType string) string {
	if mediaType == "tv" {
		return "/catalog/series/" + strconv.Itoa(id)
	}
	return "/catalog/movies/" + strconv.Itoa(id)
}

// Ajoute ou retire un film/serie des favoris du membre.
func (c *CatalogController) ToggleFavori(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, mediaType, ok := mediaParams(r)
	if !ok {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	_ = r.ParseForm()
	titre := r.FormValue("titre")
	affiche := r.FormValue("affiche")
	if err := c.filmService.ToggleFavori(user.Id, id, mediaType, titre, affiche); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, filmPath(id, mediaType)+"#social", http.StatusSeeOther)
}

// Publie un commentaire sous un film/serie.
func (c *CatalogController) AddComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, mediaType, ok := mediaParams(r)
	if !ok {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	_ = r.ParseForm()
	if err := c.filmService.AddComment(id, mediaType, user.Id, r.FormValue("contenu")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, filmPath(id, mediaType)+"#social", http.StatusSeeOther)
}

// Supprime un commentaire (auteur ou administrateur).
func (c *CatalogController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	commentId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	if err := c.filmService.DeleteComment(commentId, user.Id, user.IsAdmin()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = r.ParseForm()
	http.Redirect(w, r, r.FormValue("retour"), http.StatusSeeOther)
}

// Affiche la liste des films/series favoris du membre.
func (c *CatalogController) Favoris(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	favoris, err := c.filmService.ListFavoris(user.Id)
	if err != nil {
		http.Error(w, "Erreur lors du chargement des favoris", http.StatusInternalServerError)
		return
	}
	data := baseData(r)
	data["Favoris"] = favoris
	c.renderer.Render(w, http.StatusOK, "catalog_favoris.html", data)
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
	c.attachFilmSocial(r, data, id, "movie")
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
	c.attachFilmSocial(r, data, id, "tv")
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
