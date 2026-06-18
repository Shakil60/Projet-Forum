package controllers

// Fonctions utilitaires partagees par les controleurs (pagination, cookies, donnees de base).

import (
	"forum/middleware"
	"net/http"
	"strconv"
)

// Lit le numero de page dans l'URL, 1 par defaut.
func parsePage(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// Lit la taille de page demandee, 10 par defaut.
func parseSize(r *http.Request) (int, string) {
	switch r.URL.Query().Get("size") {
	case "all":
		return -1, "all"
	case "20":
		return 20, "20"
	case "30":
		return 30, "30"
	default:
		return 10, "10"
	}
}

// Lit l'ordre de tri demande, recent par defaut.
func parseSort(r *http.Request) string {
	switch r.URL.Query().Get("sort") {
	case "ancien", "populaire":
		return r.URL.Query().Get("sort")
	default:
		return "recent"
	}
}

// Prepare les donnees communes a tous les templates (utilisateur courant).
func baseData(r *http.Request) map[string]any {
	return map[string]any{
		"CurrentUser": middleware.GetUser(r),
	}
}

// Depose le cookie de session avec le jeton d'authentification.
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60,
	})
}

// Supprime le cookie de session.
func clearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
