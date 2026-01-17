package main

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"groupie-tracker/api"
	"groupie-tracker/models"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(400, 700))

	artists, err1 := api.GetArtists()
	relationsData, err2 := api.GetRelations()
	locations, err3 := api.GetLocations()

	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("Erreur critique : Impossible de charger les données API.")
	}

	artistesAffiches := artists

	barreRecherche := widget.NewEntry()
	barreRecherche.SetPlaceHolder("Rechercher...")

	barreRecherche.OnChanged = func(texte string) {
		artistesAffiches = RechercheArtiste(texte, artists, locations)
		list.Refresh()
	}

	list := widget.NewList(
		func() int { return len(artistesAffiches) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Nom de l'artiste")
			label.TextStyle = fyne.TextStyle{Bold: true}
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(artistesAffiches[i].Name)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		art := artistesAffiches[id]

		var rel models.Relation
		for _, r := range relationsData.Index {
			if r.ID == art.ID {
				rel = r
				break
			}
		}

		details := "Artiste : " + art.Name + "\n"
		details += "Création : " + strconv.Itoa(art.CreationDate) + "\n"
		details += "Premier album : " + art.FirstAlbum + "\n\n"

		details += "Concerts et dates:\n"
		details += "----------------------\n"

		for ville, dates := range rel.DatesLocations {
			villePropre := strings.ToUpper(strings.ReplaceAll(ville, "-", " / "))
			villePropre = strings.ReplaceAll(villePropre, "_", " ")

			details += "📍 " + villePropre + " :\n"
			for _, d := range dates {
				details += "   ▪ " + d + "\n"
			}
			details += "\n"
		}

		labelDetails := widget.NewLabel(details)
		labelDetails.Wrapping = fyne.TextWrapWord

		scrollContainer := container.NewVScroll(labelDetails)
		scrollContainer.SetMinSize(fyne.NewSize(300, 400))

		var popup *widget.PopUp
		closeButton := widget.NewButton("Fermer", func() {
			popup.Hide()
			list.Unselect(id)
		})

		content := container.NewBorder(
			nil, closeButton, nil, nil,
			scrollContainer,
		)

		popup = widget.NewModalPopUp(content, myWindow.Canvas())
		popup.Show()
	}

	contenu := container.NewBorder(barreRecherche, nil, nil, nil, list)
	myWindow.SetContent(contenu)
	myWindow.ShowAndRun()
}
