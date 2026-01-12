package api

import (
    "encoding/json"
    "net/http"
    "groupie-tracker/models" 
)

const UrlArtists = "https://groupietrackers.herokuapp.com/api/artists"
const UrlLocations = "https://groupietrackers.herokuapp.com/api/locations"

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

func GetLocations() ([]models.Location, error) {

	resp, err := http.Get(UrlLocations)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locationResponse models.LocationResponse
	err = json.NewDecoder(resp.Body).Decode(&locationResponse)
	if err != nil {
		return nil, err
	}
	return locationResponse.Index, nil
}
	

