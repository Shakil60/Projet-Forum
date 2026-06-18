package controllers

// Gere les requetes HTTP liees aux messages d'un fil de discussion.

import (
	"forum/helper"
	"forum/middleware"
	"forum/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MessageController struct {
	messageService *services.MessageService
	renderer       *helper.Renderer
}

func InitMessageController(messageService *services.MessageService, renderer *helper.Renderer) *MessageController {
	return &MessageController{messageService: messageService, renderer: renderer}
}

func readMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

// Ajoute un nouveau message dans un fil.
func (c *MessageController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	threadId, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	contenu := r.FormValue("contenu")
	if _, err := c.messageService.Create(user.Id, threadId, contenu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadId), http.StatusSeeOther)
}

// Affiche le formulaire d'edition d'un message de l'utilisateur.
func (c *MessageController) EditForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readMessageId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	message, err := c.messageService.GetForEdit(user.Id, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	data := baseData(r)
	data["Message"] = message
	c.renderer.Render(w, http.StatusOK, "message_form.html", data)
}

// Enregistre les modifications d'un message.
func (c *MessageController) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readMessageId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	contenu := r.FormValue("contenu")
	threadId, err := c.messageService.Update(user.Id, id, contenu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadId), http.StatusSeeOther)
}

// Supprime un message (auteur ou admin).
func (c *MessageController) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readMessageId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	threadId, err := c.messageService.Delete(user.Id, user.IsAdmin(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadId), http.StatusSeeOther)
}
