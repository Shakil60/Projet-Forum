package middleware

// Middlewares d'authentification : verifie le cookie et controle les acces.

import (
	"context"
	"forum/auth"
	"forum/helper"
	"forum/models"
	"forum/repositories"
	"net/http"
)

const CookieName = "token"

type contextKey string

const userContextKey contextKey = "currentUser"

type Middleware struct {
	userRepository *repositories.UserRepository
}

func InitMiddleware(userRepository *repositories.UserRepository) *Middleware {
	return &Middleware{userRepository: userRepository}
}

// Lit le cookie, valide le jeton et renvoie l'utilisateur connecte (ou nil).
func (m *Middleware) resolveUser(r *http.Request) *models.User {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return nil
	}

	claims, claimsErr := auth.ValidateToken(cookie.Value)
	if claimsErr != nil {
		return nil
	}

	user, userErr := m.userRepository.FindById(claims.UserID)
	if userErr != nil || user.Id == 0 || user.Banned {
		return nil
	}

	return &user
}

func (m *Middleware) withUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Ajoute l'utilisateur au contexte s'il est connecte, sans bloquer l'acces.
func (m *Middleware) Optional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.resolveUser(r)
		if user != nil {
			r = m.withUser(r, user)
		}
		next.ServeHTTP(w, r)
	})
}

// Redirige vers /login si l'utilisateur n'est pas connecte.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.resolveUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, m.withUser(r, user))
	})
}

// Autorise l'acces uniquement aux administrateurs connectes.
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.resolveUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !user.IsAdmin() {
			http.Error(w, "Acces reserve aux administrateurs", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, m.withUser(r, user))
	})
}

// Comme RequireAuth mais renvoie une erreur JSON pour les appels API.
func (m *Middleware) RequireAuthAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.resolveUser(r)
		if user == nil {
			helper.WriteError(w, http.StatusUnauthorized, "authentification requise")
			return
		}
		next.ServeHTTP(w, m.withUser(r, user))
	})
}

// Recupere l'utilisateur stocke dans le contexte de la requete.
func GetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
