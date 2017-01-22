package main

import (
	"github.com/zballs/go_resonate/api"
	. "github.com/zballs/go_resonate/util"
	"net/http"
)

const dir = ""

func main() {

	CreatePages("artist", "listener", "login")
	RegisterTemplates("artist.html", "login.html")

	// Create request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", TemplateHandler("artist.html"))
	mux.HandleFunc("/login", TemplateHandler("login.html"))

	// Create api
	api := api.NewApi(dir)

	// Add routes to multiplexer
	api.AddRoutes(mux)

	// Start HTTP server with multiplexer
	http.ListenAndServe(":8888", mux)
}
