package controllers

import (
	"encoding/json"
	"forum/helper"
	"forum/middleware"
	"forum/services"
	"net/http"
)

type ReactionController struct {
	reactionService *services.ReactionService
}

func InitReactionController(reactionService *services.ReactionService) *ReactionController {
	return &ReactionController{reactionService: reactionService}
}

type reactionRequest struct {
	Type string `json:"type"`
}

func (c *ReactionController) Toggle(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	messageId, idErr := readMessageId(r)
	if idErr != nil {
		helper.WriteError(w, http.StatusBadRequest, "identifiant message invalide")
		return
	}

	var body reactionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "requete invalide")
		return
	}

	result, err := c.reactionService.Toggle(user.Id, messageId, body.Type)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, result)
}
