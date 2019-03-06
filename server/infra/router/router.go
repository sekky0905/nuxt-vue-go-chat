package router

import (
	"github.com/gorilla/mux"
)

// Router is the gorilla router for API.
var Router *mux.Router

func init() {
	// Instantiation for  gorilla Router.
	r := mux.NewRouter()

	Router = r
}
