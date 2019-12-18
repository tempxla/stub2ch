package main

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/handle"
	"log"
	"net/http"
	"os"
)

func main() {

	router := handle.NewBoardRouter(nil)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
