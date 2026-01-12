package main

import (
	"fmt"
	"html/template"
	"net/http"

	// Assure-toi que le nom ici correspond bien à ton go.mod
	"groupie-tracker/api"
	"groupie-tracker/models"
)

// PageData est la "valise" qui contient toutes les données à envoyer au HTML
type PageData struct {
	Artists   []models.Artist
	Dates     models.DateList
	Relations models.RelationList
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		fmt.Println("Erreur Artistes:", err)
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}

	dates, err := api.GetDate()
	if err != nil {
		fmt.Println("Erreur Dates:", err)
		http.Error(w, "Erreur lors de la récupération des dates", http.StatusInternalServerError)
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		fmt.Println("Erreur Relations:", err)
		http.Error(w, "Erreur lors de la récupération des relations", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Artists:   artists,
		Dates:     dates,
		Relations: relations,
	}

	
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Impossible de charger le fichier index.html", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

func main() {
	fmt.Println("Le serveur démarre sur http://localhost:8080")

	
	http.HandleFunc("/", HomeHandler)

	http.ListenAndServe(":8080", nil)
}