package main

import (
	"fmt"
	"net/http"
	"groupie-tracker/api"
	"groupie-tracker/models"
)

type PageData struct {
	Artists   []models.Artist
	Dates     models.DateList
	Relations models.RelationList
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "Erreur Artistes", 500)
		return
	}

	dates, err := api.GetDate()
	if err != nil {
		http.Error(w, "Erreur Dates", 500)
		return
	}

	relations, err := api.GetRelations()
	if err != nil {
		http.Error(w, "Erreur Relations", 500)
		return
	}

	
	data := PageData{
		Artists:   artists,
		Dates:     dates,
		Relations: relations,
	}

	// (Obligatoire pour que ça compile)
	// Au lieu d'afficher sur le site, on l'affiche juste dans ton terminal
	// pour prouver que la variable 'data' est bien remplie.
	fmt.Println("Données chargées avec succès en mémoire (Artistes, Dates, Relations)")
	
	// Cette ligne sert juste à "toucher" la variable data pour que Go soit content
	_ = data 
}

func main() {
	fmt.Println("Serveur en attente sur http://localhost:8080")
	
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":8080", nil)
}