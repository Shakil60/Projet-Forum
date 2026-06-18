package routers

// Declare les routes du catalogue de films, series et personnes.

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Enregistre les routes du catalogue sur le routeur.
func RegisterCatalogRoutes(r *mux.Router, c *controllers.CatalogController, mw *middleware.Middleware) {
	r.Handle("/catalog", mw.Optional(http.HandlerFunc(c.Home))).Methods("GET")
	r.Handle("/catalog/search", mw.Optional(http.HandlerFunc(c.Search))).Methods("GET")
	r.Handle("/catalog/movies", mw.Optional(http.HandlerFunc(c.Movies))).Methods("GET")
	r.Handle("/catalog/movies/{id:[0-9]+}", mw.Optional(http.HandlerFunc(c.ShowMovie))).Methods("GET")
	r.Handle("/catalog/series", mw.Optional(http.HandlerFunc(c.Series))).Methods("GET")
	r.Handle("/catalog/series/{id:[0-9]+}", mw.Optional(http.HandlerFunc(c.ShowSeries))).Methods("GET")
	r.Handle("/catalog/people", mw.Optional(http.HandlerFunc(c.People))).Methods("GET")
	r.Handle("/catalog/people/{id:[0-9]+}", mw.Optional(http.HandlerFunc(c.ShowPerson))).Methods("GET")

	// Favoris et commentaires de films (reserves aux membres connectes)
	r.Handle("/favoris", mw.RequireAuth(http.HandlerFunc(c.Favoris))).Methods("GET")
	r.Handle("/catalog/movies/{id:[0-9]+}/favori", mw.RequireAuth(http.HandlerFunc(c.ToggleFavori))).Methods("POST")
	r.Handle("/catalog/movies/{id:[0-9]+}/commentaires", mw.RequireAuth(http.HandlerFunc(c.AddComment))).Methods("POST")
	r.Handle("/catalog/series/{id:[0-9]+}/favori", mw.RequireAuth(http.HandlerFunc(c.ToggleFavori))).Methods("POST")
	r.Handle("/catalog/series/{id:[0-9]+}/commentaires", mw.RequireAuth(http.HandlerFunc(c.AddComment))).Methods("POST")
	r.Handle("/catalog/commentaires/{id:[0-9]+}/delete", mw.RequireAuth(http.HandlerFunc(c.DeleteComment))).Methods("POST")
}
