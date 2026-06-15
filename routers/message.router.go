package routers

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterMessageRoutes(r *mux.Router, c *controllers.MessageController, mw *middleware.Middleware) {
	r.Handle("/threads/{id:[0-9]+}/messages", mw.RequireAuth(http.HandlerFunc(c.Create))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/edit", mw.RequireAuth(http.HandlerFunc(c.EditForm))).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}/edit", mw.RequireAuth(http.HandlerFunc(c.Update))).Methods("POST")
	r.Handle("/messages/{id:[0-9]+}/delete", mw.RequireAuth(http.HandlerFunc(c.Delete))).Methods("POST")
}
