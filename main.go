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
	myWindow.Resize(fyne.NewSize(500, 700))

	artists, err1 := api.GetArtists()
	relationsData, err2 := api.GetRelations()
	locations, err3 := api.GetLocations()

	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("Erreur critique : Impossible de charger les données API.")
	}

	artistesAffiches := artists

	labelDateCreation := widget.NewLabel("Date de création:")
	sliderDateMin := widget.NewSlider(1960, 2020)
	sliderDateMin.SetValue(1960)
	sliderDateMax := widget.NewSlider(1960, 2020)
	sliderDateMax.SetValue(2020)
	labelDateValues := widget.NewLabel("1960 - 2020")

	labelAlbum := widget.NewLabel("Premier album:")
	sliderAlbumMin := widget.NewSlider(1960, 2020)
	sliderAlbumMin.SetValue(1960)
	sliderAlbumMax := widget.NewSlider(1960, 2020)
	sliderAlbumMax.SetValue(2020)
	labelAlbumValues := widget.NewLabel("1960 - 2020")

	labelMembres := widget.NewLabel("Nombre de membres:")
	sliderMembresMin := widget.NewSlider(1, 10)
	sliderMembresMin.SetValue(1)
	sliderMembresMax := widget.NewSlider(1, 10)
	sliderMembresMax.SetValue(10)
	labelMembresValues := widget.NewLabel("1 - 10")

	labelLieu := widget.NewLabel("Lieu de concert:")
	lieux := RecupererLieux(locations)
	lieux = append([]string{"Tous"}, lieux...)
	selectLieu := widget.NewSelect(lieux, nil)
	selectLieu.SetSelected("Tous")

	boutonFiltrer := widget.NewButton("Appliquer les filtres", nil)

	barreRecherche := widget.NewEntry()
	barreRecherche.SetPlaceHolder("Rechercher...")

	var list *widget.List

	appliquerFiltres := func() {
		dateMin := int(sliderDateMin.Value)
		dateMax := int(sliderDateMax.Value)
		albumMin := int(sliderAlbumMin.Value)
		albumMax := int(sliderAlbumMax.Value)
		membresMin := int(sliderMembresMin.Value)
		membresMax := int(sliderMembresMax.Value)

		lieuChoisi := ""
		if selectLieu.Selected != "Tous" {
			lieuChoisi = selectLieu.Selected
		}

		artistesFiltres := FiltrerArtistes(artists, locations, dateMin, dateMax, albumMin, albumMax, membresMin, membresMax, lieuChoisi)

		if barreRecherche.Text != "" {
			artistesAffiches = RechercheArtiste(barreRecherche.Text, artistesFiltres, locations)
		} else {
			artistesAffiches = artistesFiltres
		}

		list.Refresh()
	}

	sliderDateMin.OnChanged = func(v float64) {
		labelDateValues.SetText(strconv.Itoa(int(sliderDateMin.Value)) + " - " + strconv.Itoa(int(sliderDateMax.Value)))
	}
	sliderDateMax.OnChanged = func(v float64) {
		labelDateValues.SetText(strconv.Itoa(int(sliderDateMin.Value)) + " - " + strconv.Itoa(int(sliderDateMax.Value)))
	}
	sliderAlbumMin.OnChanged = func(v float64) {
		labelAlbumValues.SetText(strconv.Itoa(int(sliderAlbumMin.Value)) + " - " + strconv.Itoa(int(sliderAlbumMax.Value)))
	}
	sliderAlbumMax.OnChanged = func(v float64) {
		labelAlbumValues.SetText(strconv.Itoa(int(sliderAlbumMin.Value)) + " - " + strconv.Itoa(int(sliderAlbumMax.Value)))
	}
	sliderMembresMin.OnChanged = func(v float64) {
		labelMembresValues.SetText(strconv.Itoa(int(sliderMembresMin.Value)) + " - " + strconv.Itoa(int(sliderMembresMax.Value)))
	}
	sliderMembresMax.OnChanged = func(v float64) {
		labelMembresValues.SetText(strconv.Itoa(int(sliderMembresMin.Value)) + " - " + strconv.Itoa(int(sliderMembresMax.Value)))
	}

	boutonFiltrer.OnTapped = appliquerFiltres

	barreRecherche.OnChanged = func(texte string) {
		appliquerFiltres()
	}

	list = widget.NewList(
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

	filtresBox := container.NewVBox(
		labelDateCreation,
		sliderDateMin,
		sliderDateMax,
		labelDateValues,
		widget.NewSeparator(),
		labelAlbum,
		sliderAlbumMin,
		sliderAlbumMax,
		labelAlbumValues,
		widget.NewSeparator(),
		labelMembres,
		sliderMembresMin,
		sliderMembresMax,
		labelMembresValues,
		widget.NewSeparator(),
		labelLieu,
		selectLieu,
		boutonFiltrer,
	)

	filtresScroll := container.NewVScroll(filtresBox)
	filtresScroll.SetMinSize(fyne.NewSize(200, 0))

	contenuPrincipal := container.NewBorder(barreRecherche, nil, nil, nil, list)

	contenuFinal := container.NewBorder(nil, nil, filtresScroll, nil, contenuPrincipal)

	myWindow.SetContent(contenuFinal)

	myWindow.ShowAndRun()
}
