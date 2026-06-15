package routers

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterReactionRoutes(r *mux.Router, c *controllers.ReactionController, mw *middleware.Middleware) {
	r.Handle("/api/messages/{id:[0-9]+}/reactions", mw.RequireAuthAPI(http.HandlerFunc(c.Toggle))).Methods("POST")
}
