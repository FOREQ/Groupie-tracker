package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	groupieapi "tracker/internal/api"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, 404, "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		HandleError(w, 405, "Method not allowed")
		return
	}
	artists, err := groupieapi.IndexArtists()
	if err != nil {
		log.Println(err)
		HandleError(w, 500, "Internal server error")
		return
	}
	log.Printf("Fetched %d artists", len(artists)) // Логируем количество артистов

	err = mainPageTemplate.Execute(w, artists)
	if err != nil {
		log.Println("Error executing template:", err)
		HandleError(w, 500, "Internal server error")
	}
	mainPageTemplate.Execute(w, artists)
}

func HandleArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, 405, "Method not allowed")
		return
	}
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if !(len(parts) == 2 && parts[0] == "artist") {
		HandleError(w, 404, "Not found")
		return
	}
	idString := parts[1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		HandleError(w, 404, "Artist not found")
		return
	}
	data, errorChannel := groupieapi.BundleArtistData(id)
	if groupieapi.ArtistNotFound(data.Artist) {
		HandleError(w, 404, "Artist not found")
		return
	}
	for len(errorChannel) > 0 {
		err = <-errorChannel
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		HandleError(w, 500, "Internal server error")
		return
	}
	if err = artistTemplate.Execute(w, data); err != nil {
		log.Println(err)
		HandleError(w, 500, "Internal server error")
		return
	}
}

func HandleError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	err := errorPageTemplate.Execute(w, struct {
		ErrorCode    int
		ErrorMessage string
	}{
		code, msg,
	})
	if err != nil {
		log.Println(err)
	}
}

const (
	errorTemplatePath  = "templates/error.html"
	indexTemplatePath  = "templates/index.html"
	artistTemplatePath = "templates/artist.html"
)

var (
	errorPageTemplate = template.Must(template.ParseFiles(errorTemplatePath))
	mainPageTemplate  = template.Must(template.ParseFiles(indexTemplatePath))
	artistTemplate    = template.Must(template.ParseFiles(artistTemplatePath))
)
