package controllers

import (
	"forum/middleware"
	"net/http"
	"strconv"
)

func parsePage(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

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

func parseSort(r *http.Request) string {
	switch r.URL.Query().Get("sort") {
	case "ancien", "populaire":
		return r.URL.Query().Get("sort")
	default:
		return "recent"
	}
}

func baseData(r *http.Request) map[string]any {
	return map[string]any{
		"CurrentUser": middleware.GetUser(r),
	}
}

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
