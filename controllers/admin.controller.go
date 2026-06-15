package controllers

import (
	"forum/helper"
	"forum/middleware"
	"forum/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AdminController struct {
	adminService  *services.AdminService
	threadService *services.ThreadService
	renderer      *helper.Renderer
}

func InitAdminController(adminService *services.AdminService, threadService *services.ThreadService, renderer *helper.Renderer) *AdminController {
	return &AdminController{
		adminService:  adminService,
		threadService: threadService,
		renderer:      renderer,
	}
}

func (c *AdminController) Dashboard(w http.ResponseWriter, r *http.Request) {
	users, usersErr := c.adminService.ListUsers()
	if usersErr != nil {
		http.Error(w, usersErr.Error(), http.StatusInternalServerError)
		return
	}

	threads, threadsErr := c.threadService.ReadAllForAdmin()
	if threadsErr != nil {
		http.Error(w, threadsErr.Error(), http.StatusInternalServerError)
		return
	}

	data := baseData(r)
	data["Users"] = users
	data["Threads"] = threads
	data["NbUsers"] = len(users)
	data["NbThreads"] = len(threads)

	c.renderer.Render(w, http.StatusOK, "admin.html", data)
}

func (c *AdminController) Ban(w http.ResponseWriter, r *http.Request) {
	c.setBan(w, r, true)
}

func (c *AdminController) Unban(w http.ResponseWriter, r *http.Request) {
	c.setBan(w, r, false)
}

func (c *AdminController) setBan(w http.ResponseWriter, r *http.Request, banned bool) {
	user := middleware.GetUser(r)

	targetId, idErr := strconv.Atoi(mux.Vars(r)["id"])
	if idErr != nil {
		http.NotFound(w, r)
		return
	}

	if err := c.adminService.SetBanned(user.Id, targetId, banned); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
