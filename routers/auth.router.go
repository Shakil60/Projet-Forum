package routers

// Declare les routes d'inscription, de connexion et de deconnexion.

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Enregistre les routes d'authentification sur le routeur.
func RegisterAuthRoutes(r *mux.Router, c *controllers.AuthController, mw *middleware.Middleware) {
	r.Handle("/register", mw.Optional(http.HandlerFunc(c.RegisterForm))).Methods("GET")
	r.Handle("/register", mw.Optional(http.HandlerFunc(c.Register))).Methods("POST")
	r.Handle("/login", mw.Optional(http.HandlerFunc(c.LoginForm))).Methods("GET")
	r.Handle("/login", mw.Optional(http.HandlerFunc(c.Login))).Methods("POST")
	r.Handle("/logout", http.HandlerFunc(c.Logout)).Methods("POST")
}
