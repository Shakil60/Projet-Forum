package dto

// Decrit le format JSON renvoye en cas d'erreur de l'API.

type ApiError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
