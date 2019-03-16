package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
	"github.com/sekky0905/nuxt-vue-go-chat/server/interface/controller"
)

func main() {
	router.G.GET("/", func(g *gin.Context) {
		g.File("../client/nuxt-vue-go-chat/dist/index.html")
	})

	apiV1 := router.G.Group("/v1")

	dbm := db.NewDBManager()
	ac := initializeAuthenticationController(dbm)
	ac.InitAuthenticationAPI(apiV1)

	if err := router.G.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

// initializeAuthenticationController generates and returns AuthenticationController.
func initializeAuthenticationController(m repository.DBManager) controller.AuthenticationController {
	txCloser := db.CloseTransaction

	uRepo := db.NewUserRepository()
	sRepo := db.NewSessionRepository()
	uService := service.NewUserService(m, uRepo)
	sService := service.NewSessionService(m, sRepo)

	di := application.NewAuthenticationServiceDIInput(uRepo, sRepo, uService, sService)
	aApp := application.NewAuthenticationService(m, di, txCloser)

	return controller.NewAuthenticationController(aApp)
}
