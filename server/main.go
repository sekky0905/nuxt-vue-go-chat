package main

import (
	"net/http"

	"github.com/sekky0905/nuxt-vue-go-chat/server/infra/router"
)

func main() {
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
