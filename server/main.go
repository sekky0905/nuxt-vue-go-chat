package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db/query"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
	"github.com/sekky0905/nuxt-vue-go-chat/server/interface/controller"
	"github.com/sekky0905/nuxt-vue-go-chat/server/middleware"
)

func main() {
	apiV1 := router.G.Group("/v1")

	dbm := db.NewDBManager()
	ac := initializeAuthenticationController(dbm)
	ac.InitAuthenticationAPI(apiV1)

	threadRouting := apiV1.Group("/threads")

	// use middleware
	threadRouting.Use(middleware.CheckAuthentication())

	cc := initializeCommentController(dbm)
	cc.InitCommentAPI(threadRouting)

	tc := initializeThreadController(dbm)
	tc.InitThreadAPI(threadRouting)

	router.G.NoRoute(func(g *gin.Context) {
		g.File("./../client/nuxt-vue-go-chat/dist/index.html")
	})
	router.G.Static("/_nuxt", "./../client/nuxt-vue-go-chat/dist/_nuxt/")

	if err := router.G.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

// initializeAuthenticationController generates and returns AuthenticationController.
func initializeAuthenticationController(m query.DBManager) controller.AuthenticationController {
	txCloser := db.CloseTransaction

	uRepo := db.NewUserRepository()
	sRepo := db.NewSessionRepository()
	uService := service.NewUserService(m, uRepo)
	sService := service.NewSessionService(sRepo)
	aService := service.NewAuthenticationService(uRepo)

	di := application.NewAuthenticationServiceDIInput(uRepo, sRepo, uService, sService, aService)
	aApp := application.NewAuthenticationService(m, di, txCloser)

	return controller.NewAuthenticationController(aApp)
}

// initializeThreadCController generates and returns ThreadCController.
func initializeThreadController(m query.DBManager) controller.ThreadController {
	txCloser := db.CloseTransaction

	tRepo := db.NewThreadRepository()
	tService := service.NewThreadService(tRepo)

	tApp := application.NewThreadService(m, tService, tRepo, txCloser)

	return controller.NewThreadController(tApp)
}

// initializeCommentController generates and returns CommentController.
func initializeCommentController(m query.DBManager) controller.CommentController {
	txCloser := db.CloseTransaction

	cRepo := db.NewCommentRepository()
	cService := service.NewCommentService(cRepo)

	cApp := application.NewCommentService(m, cService, cRepo, txCloser)

	return controller.NewCommentController(cApp)
}
