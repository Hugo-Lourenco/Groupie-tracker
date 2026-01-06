package api

import (
    "encoding/json"
    "net/http"
    "groupie-tracker/models" 
)

const UrlArtists = "https://groupietrackers.herokuapp.com/api/artists"

func GetArtists() ([]models.Artist, error ) {

	resp, err := http.Get(UrlArtists)

	if err != nil { 
		return nil, err 
	}
	defer resp.Body.Close()	

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
		if err != nil { 
		return nil, err 
	}
	return artists, nil
}
	

