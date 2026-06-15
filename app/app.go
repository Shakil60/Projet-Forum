package app

import (
	"database/sql"
	"net/http"

	"forum/auth"
	"forum/config"
	"forum/controllers"
	"forum/helper"
	"forum/middleware"
	"forum/repositories"
	"forum/routers"
	"forum/services"

	"github.com/gorilla/mux"
)

type App struct {
	Db     *sql.DB
	Router *mux.Router
	Port   string
}

func InitApp() *App {
	config.LoadEnv()
	auth.SetSecret(config.GetEnvWithDefault("JWT_SECRET", "cinetalk_dev_secret_change_me"))

	db := config.InitDB()
	renderer := helper.InitRenderer("./views")

	userRepository := repositories.InitUserRepository(db)
	tagRepository := repositories.InitTagRepository(db)
	threadRepository := repositories.InitThreadRepository(db)
	messageRepository := repositories.InitMessageRepository(db)
	reactionRepository := repositories.InitReactionRepository(db)

	authService := services.InitAuthService(userRepository)
	tagService := services.InitTagService(tagRepository)
	threadService := services.InitThreadService(threadRepository, tagRepository)
	messageService := services.InitMessageService(messageRepository, threadRepository)
	reactionService := services.InitReactionService(reactionRepository, messageRepository)
	adminService := services.InitAdminService(userRepository)

	authController := controllers.InitAuthController(authService, renderer)
	threadController := controllers.InitThreadController(threadService, messageService, tagService, renderer)
	messageController := controllers.InitMessageController(messageService, renderer)
	reactionController := controllers.InitReactionController(reactionService)
	adminController := controllers.InitAdminController(adminService, threadService, renderer)

	mw := middleware.InitMiddleware(userRepository)

	router := mux.NewRouter()

	routers.RegisterThreadRoutes(router, threadController, mw)
	routers.RegisterAuthRoutes(router, authController, mw)
	routers.RegisterMessageRoutes(router, messageController, mw)
	routers.RegisterReactionRoutes(router, reactionController, mw)
	routers.RegisterAdminRoutes(router, adminController, mw)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return &App{
		Db:     db,
		Router: router,
		Port:   config.GetEnvWithDefault("SERVER_PORT", "8080"),
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
