package main

import (
	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
	"github.com/sekky0905/nuxt-vue-go-chat/server/interface/controller"
)

func main() {
	router.G.Static("/", "../client/nuxt-vue-go-chat/dist")
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
	aService := service.NewAuthenticationService(m, uRepo)

	di := application.NewAuthenticationServiceDIInput(uRepo, sRepo, uService, sService, aService)
	aApp := application.NewAuthenticationService(m, di, txCloser)

	return controller.NewAuthenticationController(aApp)
}
