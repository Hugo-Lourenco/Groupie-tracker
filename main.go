package main

import (
    "fmt"
    "net/http"
    "groupie-tracker/api"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    artists, err := api.GetArtists()
    
    if err != nil {
        http.Error(w, "Erreur du serveur", http.StatusInternalServerError)
        return
    }

   
    fmt.Fprintf(w, "J'ai trouvé %d artistes !", len(artists))
}

func main() {
    fmt.Println("Le serveur démarre sur http://localhost:8080")
    
    http.HandleFunc("/", HomeHandler)
    http.ListenAndServe(":8080", nil)
}