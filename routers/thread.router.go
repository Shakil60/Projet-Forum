package routers

// Declare les routes liees aux fils de discussion.

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Enregistre les routes des fils sur le routeur.
func RegisterThreadRoutes(r *mux.Router, c *controllers.ThreadController, mw *middleware.Middleware) {
	r.Handle("/", mw.Optional(http.HandlerFunc(c.Home))).Methods("GET")

	r.Handle("/threads/new", mw.RequireAuth(http.HandlerFunc(c.NewForm))).Methods("GET")
	r.Handle("/threads", mw.RequireAuth(http.HandlerFunc(c.Create))).Methods("POST")

	r.Handle("/threads/{id:[0-9]+}", mw.Optional(http.HandlerFunc(c.Show))).Methods("GET")
	r.Handle("/threads/{id:[0-9]+}/edit", mw.RequireAuth(http.HandlerFunc(c.EditForm))).Methods("GET")
	r.Handle("/threads/{id:[0-9]+}/edit", mw.RequireAuth(http.HandlerFunc(c.Update))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/state", mw.RequireAuth(http.HandlerFunc(c.ChangeState))).Methods("POST")
	r.Handle("/threads/{id:[0-9]+}/delete", mw.RequireAuth(http.HandlerFunc(c.Delete))).Methods("POST")
}
