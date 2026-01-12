package api

import (
    "encoding/json"
    "net/http"
    "groupie-tracker/models" 
)

const UrlArtists	 = 		"https://groupietrackers.herokuapp.com/api/artists"
const UrlLocations	 = 		"https://groupietrackers.herokuapp.com/api/locations"
const UrlDate 		 = 		"https://groupietrackers.herokuapp.com/api/dates"

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
	

func GetDate() (models.DateList, error) {

	resp, err := http.Get(UrlDate)
	if err != nil {
		return models.DateList{}, err
	}
	defer resp.Body.Close()

	var data models.DateList
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return models.DateList{}, err
	}
	return data, nil 
}
