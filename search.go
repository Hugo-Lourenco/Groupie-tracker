package main

import (
	"groupie-tracker/models"
	"strconv"
	"strings"
)

func RechercheArtiste(recherche string, artistes []models.Artist, locations []models.Location) []models.Artist {
	recherche = strings.ToLower(recherche)

	if recherche == "" {
		return artistes
	}

	resultat := []models.Artist{}

	for _, artiste := range artistes {
		trouve := false

		if strings.Contains(strings.ToLower(artiste.Name), recherche) {
			trouve = true
		}

		for _, membre := range artiste.Members {
			if strings.Contains(strings.ToLower(membre), recherche) {
				trouve = true
			}
		}

		if strings.Contains(strconv.Itoa(artiste.CreationDate), recherche) {
			trouve = true
		}

		if strings.Contains(strings.ToLower(artiste.FirstAlbum), recherche) {
			trouve = true
		}

		for _, loc := range locations {
			if loc.ID == artiste.ID {
				for _, lieu := range loc.Locations {
					if strings.Contains(strings.ToLower(lieu), recherche) {
						trouve = true
					}
				}
			}
		}

		if trouve {
			resultat = append(resultat, artiste)
		}
	}

	return resultat
}
