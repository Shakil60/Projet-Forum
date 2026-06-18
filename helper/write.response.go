package helper

// Fonctions utilitaires pour ecrire les reponses HTTP au format JSON.

import (
	"encoding/json"
	"forum/dto"
	"net/http"
)

// Ecrit une reponse JSON avec le code de statut donne.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// Ecrit une erreur JSON standardisee avec son message.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, dto.ApiError{Status: status, Error: message})
}
