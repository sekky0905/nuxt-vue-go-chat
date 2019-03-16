package main

import (
	"net/http"

	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/service"

	"github.com/sekky0905/nuxt-vue-go-chat/server/application"
	"github.com/sekky0905/nuxt-vue-go-chat/server/domain/repository"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/db"
	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
	"github.com/sekky0905/nuxt-vue-go-chat/server/interface/controller"
)

func main() {
	dbManager := db.NewDBManager()

	apiRouter := router.Router.PathPrefix("/v1/").Subrouter()

	aController := initializeAuthenticationController(dbManager)
	apiRouter.HandleFunc("/signUp", aController.SignUp).Methods(http.MethodPost)

	// for static file
	entryPoint := "../client/nuxt-vue-go-chat/dist/index.html"
	router.Router.Path("/").HandlerFunc(ServeStaticFile(entryPoint))
	router.Router.PathPrefix("/_nuxt/").Handler(http.StripPrefix("/_nuxt/", http.FileServer(http.Dir("../client/nuxt-vue-go-chat/dist/_nuxt/"))))

	if err := http.ListenAndServe(":8080", router.Router); err != nil {
		panic(err.Error())
	}
}

// ServeStaticFile is deliver static files.
func ServeStaticFile(entryPoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entryPoint)
	}
	return http.HandlerFunc(fn)
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

	rm := router.NewRequestManager()
	return controller.NewAuthenticationController(rm, aApp)
}
