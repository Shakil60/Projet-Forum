package controllers

// Gere l'inscription, la connexion et la deconnexion des utilisateurs.

import (
	"forum/helper"
	"forum/services"
	"net/http"
)

type AuthController struct {
	authService *services.AuthService
	renderer    *helper.Renderer
}

func InitAuthController(authService *services.AuthService, renderer *helper.Renderer) *AuthController {
	return &AuthController{authService: authService, renderer: renderer}
}

// Affiche le formulaire d'inscription.
func (c *AuthController) RegisterForm(w http.ResponseWriter, r *http.Request) {
	c.renderer.Render(w, http.StatusOK, "register.html", baseData(r))
}

// Cree un compte puis connecte directement le nouvel utilisateur.
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.renderer.Render(w, http.StatusBadRequest, "register.html", baseData(r))
		return
	}

	username := r.FormValue("nom_utilisateur")
	email := r.FormValue("email")
	password := r.FormValue("mot_de_passe")

	user, err := c.authService.Register(username, email, password)
	if err != nil {
		data := baseData(r)
		data["Error"] = err.Error()
		data["Username"] = username
		data["Email"] = email
		c.renderer.Render(w, http.StatusBadRequest, "register.html", data)
		return
	}

	token, _, loginErr := c.authService.Login(user.Username, password)
	if loginErr != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	setAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Affiche le formulaire de connexion.
func (c *AuthController) LoginForm(w http.ResponseWriter, r *http.Request) {
	c.renderer.Render(w, http.StatusOK, "login.html", baseData(r))
}

// Verifie les identifiants et depose le cookie d'authentification.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.renderer.Render(w, http.StatusBadRequest, "login.html", baseData(r))
		return
	}

	identifiant := r.FormValue("identifiant")
	password := r.FormValue("mot_de_passe")

	token, _, err := c.authService.Login(identifiant, password)
	if err != nil {
		data := baseData(r)
		data["Error"] = err.Error()
		data["Identifiant"] = identifiant
		c.renderer.Render(w, http.StatusUnauthorized, "login.html", data)
		return
	}

	setAuthCookie(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Supprime le cookie de session et deconnecte l'utilisateur.
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
