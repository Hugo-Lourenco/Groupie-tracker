package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
	"groupie-tracker/models"
)

var templates *template.Template

type PageData struct {
	Artists   []models.Artist
	Dates     models.DateList
	Relations models.RelationList
}

type ArtistPageData struct {
	Artist   models.Artist
	Relation models.Relation
}

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := PageData{
		Artists: artists,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(path)
	if err != nil || id < 1 {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	artists, err := api.GetArtists()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var artist models.Artist
	found := false
	for _, a := range artists {
		if a.ID == id {
			artist = a
			found = true
			break
		}
	}

	if !found {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	var relation models.Relation
	for _, rel := range relations.Index {
		if rel.ID == id {
			relation = rel
			break
		}
	}

	data := ArtistPageData{
		Artist:   artist,
		Relation: relation,
	}

	err = templates.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	switch status {
	case http.StatusNotFound:
		templates.ExecuteTemplate(w, "404.html", nil)
	case http.StatusInternalServerError:
		templates.ExecuteTemplate(w, "500.html", nil)
	default:
		templates.ExecuteTemplate(w, "500.html", nil)
	}
}

func main() {
	fmt.Println("Serveur Groupie Tracker demarre sur http://localhost:8080")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist/", ArtistHandler)

	http.ListenAndServe(":8080", nil)
}
