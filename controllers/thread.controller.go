package controllers

import (
	"forum/helper"
	"forum/middleware"
	"forum/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ThreadController struct {
	threadService  *services.ThreadService
	messageService *services.MessageService
	tagService     *services.TagService
	renderer       *helper.Renderer
}

func InitThreadController(threadService *services.ThreadService, messageService *services.MessageService, tagService *services.TagService, renderer *helper.Renderer) *ThreadController {
	return &ThreadController{
		threadService:  threadService,
		messageService: messageService,
		tagService:     tagService,
		renderer:       renderer,
	}
}

func readThreadId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func splitTags(raw string) []string {
	parts := strings.Split(raw, ",")
	var tags []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			tags = append(tags, part)
		}
	}
	return tags
}

func parseTagIds(values []string) []int {
	var ids []int
	for _, value := range values {
		if id, err := strconv.Atoi(value); err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

func (c *ThreadController) Home(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	search := strings.TrimSpace(r.URL.Query().Get("q"))

	if search != "" && user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tagId, _ := strconv.Atoi(r.URL.Query().Get("tag"))
	page := parsePage(r)
	size, sizeParam := parseSize(r)

	threads, pagination, err := c.threadService.ReadForList(tagId, search, page, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tags, _ := c.tagService.ReadAll()

	activeTagName := ""
	if tagId > 0 {
		if tag, tagErr := c.tagService.FindById(tagId); tagErr == nil {
			activeTagName = tag.Nom
		}
	}

	data := baseData(r)
	data["Threads"] = threads
	data["Tags"] = tags
	data["Pagination"] = pagination
	data["ActiveTagId"] = tagId
	data["ActiveTagName"] = activeTagName
	data["Search"] = search
	data["SizeParam"] = sizeParam

	c.renderer.Render(w, http.StatusOK, "home.html", data)
}

func (c *ThreadController) Show(w http.ResponseWriter, r *http.Request) {
	id, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	thread, err := c.threadService.ReadVisibleById(id)
	if err != nil {
		c.renderer.Render(w, http.StatusNotFound, "notfound.html", baseData(r))
		return
	}

	user := middleware.GetUser(r)
	currentUserId := 0
	canManage := false
	if user != nil {
		currentUserId = user.Id
		canManage = user.Id == thread.AuthorId || user.IsAdmin()
	}

	sort := parseSort(r)
	page := parsePage(r)
	size, sizeParam := parseSize(r)

	messages, pagination, msgErr := c.messageService.ReadForThread(thread.Id, sort, currentUserId, page, size)
	if msgErr != nil {
		http.Error(w, msgErr.Error(), http.StatusInternalServerError)
		return
	}

	data := baseData(r)
	data["Thread"] = thread
	data["Messages"] = messages
	data["Pagination"] = pagination
	data["Sort"] = sort
	data["SizeParam"] = sizeParam
	data["CanManage"] = canManage

	c.renderer.Render(w, http.StatusOK, "thread.html", data)
}

func (c *ThreadController) NewForm(w http.ResponseWriter, r *http.Request) {
	tags, _ := c.tagService.ReadAll()
	data := baseData(r)
	data["Tags"] = tags
	data["SelectedTags"] = map[int]bool{}
	c.renderer.Render(w, http.StatusOK, "thread_form.html", data)
}

func (c *ThreadController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	titre := r.FormValue("titre")
	contenu := r.FormValue("contenu")
	tagIds := parseTagIds(r.Form["tags"])
	newTags := splitTags(r.FormValue("nouveaux_tags"))

	threadId, err := c.threadService.Create(user.Id, titre, contenu, tagIds, newTags)
	if err != nil {
		tags, _ := c.tagService.ReadAll()
		data := baseData(r)
		data["Tags"] = tags
		data["Error"] = err.Error()
		data["Titre"] = titre
		data["Contenu"] = contenu
		data["SelectedTags"] = map[int]bool{}
		c.renderer.Render(w, http.StatusBadRequest, "thread_form.html", data)
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(threadId), http.StatusSeeOther)
}

func (c *ThreadController) EditForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	thread, err := c.threadService.GetForEdit(user.Id, user.IsAdmin(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	selected := map[int]bool{}
	for _, tag := range thread.Tags {
		selected[tag.Id] = true
	}

	tags, _ := c.tagService.ReadAll()
	data := baseData(r)
	data["Tags"] = tags
	data["SelectedTags"] = selected
	data["IsEdit"] = true
	data["ThreadId"] = thread.Id
	data["Titre"] = thread.Titre
	data["Contenu"] = thread.Contenu
	data["Etat"] = thread.Etat
	c.renderer.Render(w, http.StatusOK, "thread_form.html", data)
}

func (c *ThreadController) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	titre := r.FormValue("titre")
	contenu := r.FormValue("contenu")
	tagIds := parseTagIds(r.Form["tags"])
	newTags := splitTags(r.FormValue("nouveaux_tags"))

	if err := c.threadService.Update(user.Id, user.IsAdmin(), id, titre, contenu, tagIds, newTags); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (c *ThreadController) ChangeState(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	etat := r.FormValue("etat")
	if err := c.threadService.ChangeState(user.Id, user.IsAdmin(), id, etat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if etat == "archive" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (c *ThreadController) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, idErr := readThreadId(r)
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	if err := c.threadService.Delete(user.Id, user.IsAdmin(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
