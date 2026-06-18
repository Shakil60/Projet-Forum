package routers

// Declare la route API des reactions (like / dislike).

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Enregistre la route des reactions sur le routeur.
func RegisterReactionRoutes(r *mux.Router, c *controllers.ReactionController, mw *middleware.Middleware) {
	r.Handle("/api/messages/{id:[0-9]+}/reactions", mw.RequireAuthAPI(http.HandlerFunc(c.Toggle))).Methods("POST")
}
