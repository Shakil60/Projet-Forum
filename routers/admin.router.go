package routers

// Declare les routes du panneau d'administration.

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Enregistre les routes d'administration sur le routeur (protegees par RequireAdmin).
func RegisterAdminRoutes(r *mux.Router, c *controllers.AdminController, mw *middleware.Middleware) {
	r.Handle("/admin", mw.RequireAdmin(http.HandlerFunc(c.Dashboard))).Methods("GET")
	r.Handle("/admin/users/{id:[0-9]+}/ban", mw.RequireAdmin(http.HandlerFunc(c.Ban))).Methods("POST")
	r.Handle("/admin/users/{id:[0-9]+}/unban", mw.RequireAdmin(http.HandlerFunc(c.Unban))).Methods("POST")
}
