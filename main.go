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
)

func main() {
	// 1. Création de l'application
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(400, 700))

	// 2. Chargement des données
	// On ne charge QUE les Artistes et les Relations (Relations contient déjà les dates !)
	artists, err1 := api.GetArtists()
	relationsData, err2 := api.GetRelations()

	// Si l'un des deux échoue
	if err1 != nil || err2 != nil {
		log.Fatal("Erreur critique : Impossible de charger les données API.")
	}

	// 3. Création de la liste des noms
	list := widget.NewList(
		func() int { return len(artists) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("Nom de l'artiste")
			label.TextStyle = fyne.TextStyle{Bold: true}
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(artists[i].Name)
		},
	)

	// 4. Gestion du clic (Popup détails)
	list.OnSelected = func(id widget.ListItemID) {
		art := artists[id]
		rel := relationsData.Index[id]

		// Construction du texte
		details := "Artiste : " + art.Name + "\n"
		details += "Création : " + strconv.Itoa(art.CreationDate) + "\n"
		details += "Premier album : " + art.FirstAlbum + "\n\n"
		
		details += "Concerts et dates:\n"
		details += "----------------------\n"

		// On boucle sur la Map des relations
		for ville, dates := range rel.DatesLocations {
			// Petit nettoyage du texte pour faire joli
			villePropre := strings.ToUpper(strings.ReplaceAll(ville, "-", " / "))
			villePropre = strings.ReplaceAll(villePropre, "_", " ")

			details += "📍 " + villePropre + " :\n"
			for _, d := range dates {
				details += "   ▫ " + d + "\n"
			}
			details += "\n"
		}

		// Mise en page du popup
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

	myWindow.SetContent(container.NewMax(list))
	myWindow.ShowAndRun()
}
