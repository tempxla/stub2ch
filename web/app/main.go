package main

import (
	"./config"
	"./service"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	repo := &service.BoardStore{
		Context: ctx,
		Client:  client,
	}
	sv := service.NewBoardService(repo)

	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/:board/bbs.cgi", handleBbsCgi)
	router.GET("/:board/subject.txt", handleSubjectTxt(sv))
	router.GET("/:board/setting.txt", handleSettingTxt)
	router.GET("/:board/dat/:dat", handleDat(sv))

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
