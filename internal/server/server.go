package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"tracker/internal/handlers"
	"tracker/internal/helpers"
	"tracker/static"
)

func Init(link, port string) {
	setupErr := prepareEnv()
	if setupErr != nil {
		log.Fatal(setupErr)
	}
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/artist/", handlers.HandleArtist)

	fs := http.FileServer(http.FS(static.Files))
	staticFilesHandler := http.StripPrefix("/static/", fs)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		// path := r.URL.Path[len("/static/"):]
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		path = filepath.Join(staticFilesLocation, path)
		if path == staticFilesLocation {
			handlers.HandleError(w, 403, "Forbidden")
			return
		}
		staticFilesHandler.ServeHTTP(w, r)
	})

	addr := link + port
	log.Println("Starting server at", addr)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func prepareEnv() error {
	if err := helpers.ChangeDirProjectRoot(); err != nil {
		return err
	}
	if err := checkStaticFiles(); err != nil {
		return err
	}
	return nil
}

var staticFilesLocation = filepath.Join("static")

var staticFilenames = [...]string{
	"style.css",
	"favicon.ico",
}

func checkStaticFiles() error {
	for _, filename := range staticFilenames {
		if _, err := os.Stat(filepath.Join(staticFilesLocation, filename)); err != nil {
			return fmt.Errorf("cannot access \"%s\": %w", filename, err)
		}
	}
	return nil
}
