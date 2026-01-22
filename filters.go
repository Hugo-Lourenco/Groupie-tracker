package main

import (
	"groupie-tracker/models"
)

func FiltrerArtistes(artistes []models.Artist, locations []models.Location, dateMin, dateMax, albumMin, albumMax, nbMembresMin, nbMembresMax int, lieuChoisi string) []models.Artist {
	resultat := []models.Artist{}

	for _, artiste := range artistes {
		garde := true

		if artiste.CreationDate < dateMin || artiste.CreationDate > dateMax {
			garde = false
		}

		anneeAlbum := ExtraireAnnee(artiste.FirstAlbum)
		if anneeAlbum < albumMin || anneeAlbum > albumMax {
			garde = false
		}

		nbMembres := len(artiste.Members)
		if nbMembres < nbMembresMin || nbMembres > nbMembresMax {
			garde = false
		}

		if lieuChoisi != "" {
			aConcertIci := false
			for _, loc := range locations {
				if loc.ID == artiste.ID {
					for _, lieu := range loc.Locations {
						if lieu == lieuChoisi {
							aConcertIci = true
							break
						}
					}
				}
			}
			if !aConcertIci {
				garde = false
			}
		}

		if garde {
			resultat = append(resultat, artiste)
		}
	}

	return resultat
}

func ExtraireAnnee(date string) int {
	if len(date) >= 4 {
		annee := 0
		for i := len(date) - 4; i < len(date); i++ {
			if date[i] >= '0' && date[i] <= '9' {
				annee = annee*10 + int(date[i]-'0')
			}
		}
		return annee
	}
	return 0
}

func RecupererLieux(locations []models.Location) []string {
	lieux := []string{}
	dejaVu := make(map[string]bool)

	for _, loc := range locations {
		for _, lieu := range loc.Locations {
			if !dejaVu[lieu] {
				lieux = append(lieux, lieu)
				dejaVu[lieu] = true
			}
		}
	}

	return lieux
}
