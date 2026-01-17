package main

import (
	"strconv"
	"strings"

	"groupie-tracker/models"
)

func RechercheArtiste(motCle string, artistes []models.Artist, locations []models.Location) []models.Artist {
	motCle = strings.ToLower(motCle)

	if motCle == "" {
		return artistes
	}

	resultat := []models.Artist{}

	for _, artiste := range artistes {
		trouve := false

		if strings.Contains(strings.ToLower(artiste.Name), motCle) {
			trouve = true
		}

		if !trouve {
			for _, membre := range artiste.Members {
				if strings.Contains(strings.ToLower(membre), motCle) {
					trouve = true
					break
				}
			}
		}

		if !trouve {
			dateCreation := strconv.Itoa(artiste.CreationDate)
			if strings.Contains(dateCreation, motCle) {
				trouve = true
			}
		}

		if !trouve {
			if strings.Contains(strings.ToLower(artiste.FirstAlbum), motCle) {
				trouve = true
			}
		}

		if !trouve {
			for _, location := range locations {
				if location.ID == artiste.ID {
					for _, lieu := range location.Locations {
						if strings.Contains(strings.ToLower(lieu), motCle) {
							trouve = true
							break
						}
					}
				}
				if trouve {
					break
				}
			}
		}

		if trouve {
			resultat = append(resultat, artiste)
		}
	}

	return resultat
}
