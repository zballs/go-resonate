package main

import (
	"github.com/zballs/go_resonate/api"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

func main() {
	RegisterTemplates("home.html")
	CreatePages("home")

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/home", TemplateHandler("home.html"))

	// Create api
	api := api.NewApi()

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
