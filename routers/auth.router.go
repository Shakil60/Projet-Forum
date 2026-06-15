package routers

import (
	"forum/controllers"
	"forum/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router, c *controllers.AuthController, mw *middleware.Middleware) {
	r.Handle("/register", mw.Optional(http.HandlerFunc(c.RegisterForm))).Methods("GET")
	r.Handle("/register", mw.Optional(http.HandlerFunc(c.Register))).Methods("POST")
	r.Handle("/login", mw.Optional(http.HandlerFunc(c.LoginForm))).Methods("GET")
	r.Handle("/login", mw.Optional(http.HandlerFunc(c.Login))).Methods("POST")
	r.Handle("/logout", http.HandlerFunc(c.Logout)).Methods("POST")
}
