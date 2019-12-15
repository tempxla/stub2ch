package main

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/handle"
	"log"
	"net/http"
	"os"
)

func main() {

	router := handle.NewBoardRouter()

	// Serve static files out of the public directory.
	// By configuring a static handler in app.yaml, App Engine serves all the
	// static content itself. As a result, the following two lines are in
	// effect for development only.
	public := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	http.Handle("/static/", public)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
