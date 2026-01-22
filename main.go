package main

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"groupie-tracker/api"
	"groupie-tracker/models"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Groupie Tracker")
	myWindow.Resize(fyne.NewSize(500, 700))

	LoadLikes()

	artists, err1 := api.GetArtists()
	relationsData, err2 := api.GetRelations()
	locations, err3 := api.GetLocations()

	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("Erreur critique : Impossible de charger les données API.")
	}

	artistesAffiches := artists

	checkFavorites := widget.NewCheck(" Afficher seulement mes favoris", nil)

	labelDateCreation := widget.NewLabel("Date de création:")
	sliderDateMin := widget.NewSlider(1950, 2025)
	sliderDateMin.SetValue(1950)
	sliderDateMax := widget.NewSlider(1950, 2025)
	sliderDateMax.SetValue(2025)
	labelDateValues := widget.NewLabel("1950 - 2025")

	labelAlbum := widget.NewLabel("Premier album:")
	sliderAlbumMin := widget.NewSlider(1950, 2025)
	sliderAlbumMin.SetValue(1950)
	sliderAlbumMax := widget.NewSlider(1950, 2025)
	sliderAlbumMax.SetValue(2025)
	labelAlbumValues := widget.NewLabel("1950 - 2025")

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

	boutonFiltrer := widget.NewButtonWithIcon("Appliquer les filtres", theme.SearchIcon(), nil)
	barreRecherche := widget.NewEntry()
	barreRecherche.SetPlaceHolder("Rechercher un artiste, membre...")

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

		tempArtistes := FiltrerArtistes(artists, locations, dateMin, dateMax, albumMin, albumMax, membresMin, membresMax, lieuChoisi)

		if barreRecherche.Text != "" {
			tempArtistes = RechercheArtiste(barreRecherche.Text, tempArtistes, locations)
		}

		if checkFavorites.Checked {
			onlyLikes := []models.Artist{}
			for _, art := range tempArtistes {
				if IsLiked(art.ID) {
					onlyLikes = append(onlyLikes, art)
				}
			}
			artistesAffiches = onlyLikes
		} else {
			artistesAffiches = tempArtistes
		}

		if list != nil {
			list.Refresh()
		}
	}

	updateLabels := func() {
		labelDateValues.SetText(strconv.Itoa(int(sliderDateMin.Value)) + " - " + strconv.Itoa(int(sliderDateMax.Value)))
		labelAlbumValues.SetText(strconv.Itoa(int(sliderAlbumMin.Value)) + " - " + strconv.Itoa(int(sliderAlbumMax.Value)))
		labelMembresValues.SetText(strconv.Itoa(int(sliderMembresMin.Value)) + " - " + strconv.Itoa(int(sliderMembresMax.Value)))
	}
	sliderDateMin.OnChanged = func(f float64) { updateLabels() }
	sliderDateMax.OnChanged = func(f float64) { updateLabels() }
	sliderAlbumMin.OnChanged = func(f float64) { updateLabels() }
	sliderAlbumMax.OnChanged = func(f float64) { updateLabels() }
	sliderMembresMin.OnChanged = func(f float64) { updateLabels() }
	sliderMembresMax.OnChanged = func(f float64) { updateLabels() }

	boutonFiltrer.OnTapped = appliquerFiltres
	barreRecherche.OnChanged = func(s string) { appliquerFiltres() }
	checkFavorites.OnChanged = func(b bool) { appliquerFiltres() }

	list = widget.NewList(
		func() int { return len(artistesAffiches) },
		func() fyne.CanvasObject {
			icon := widget.NewIcon(theme.CheckButtonCheckedIcon())
			label := widget.NewLabel("Nom Artiste")
			label.TextStyle = fyne.TextStyle{Bold: true}
			return container.NewHBox(icon, label)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			icon := box.Objects[0].(*widget.Icon)
			label := box.Objects[1].(*widget.Label)

			art := artistesAffiches[i]
			label.SetText(art.Name)

			if IsLiked(art.ID) {
				icon.Show()
			} else {
				icon.Hide()
			}
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

		btnLike := widget.NewButtonWithIcon("Mettre en favoris", theme.ContentAddIcon(), nil)
		updateLikeBtn := func() {
			if IsLiked(art.ID) {
				btnLike.SetText("Retirer des favoris")
				btnLike.SetIcon(theme.DeleteIcon())
			} else {
				btnLike.SetText("Mettre en favoris")
				btnLike.SetIcon(theme.ConfirmIcon())
			}
		}
		updateLikeBtn()

		btnLike.OnTapped = func() {
			ToggleLike(art.ID)
			updateLikeBtn()
			list.Refresh()
		}

		btnSpotify := widget.NewButtonWithIcon("Ecouter sur Spotify", theme.MediaPlayIcon(), func() {
			searchURL := "https://open.spotify.com/search/" + url.QueryEscape(art.Name)
			u, _ := url.Parse(searchURL)
			myApp.OpenURL(u)
		})

		infoContainer := container.NewVBox(
			widget.NewLabelWithStyle("🎤 "+art.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(2, btnLike, btnSpotify),
			widget.NewSeparator(),
			widget.NewLabel("Création : "+strconv.Itoa(art.CreationDate)),
			widget.NewLabel("1er Album : "+art.FirstAlbum),
			widget.NewLabel("Membres : "+strings.Join(art.Members, ", ")),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Concerts & Géolocalisation", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		)

		datesContainer := container.NewVBox()

		for ville, dates := range rel.DatesLocations {
			villePropre := strings.ToUpper(strings.ReplaceAll(ville, "-", " / "))
			villePropre = strings.ReplaceAll(villePropre, "_", " ")

			btnMap := widget.NewButtonWithIcon("Voir carte", theme.SearchIcon(), func() {
				mapURL := "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(villePropre)
				u, _ := url.Parse(mapURL)
				myApp.OpenURL(u)
			})

			datesTxt := ""
			for _, d := range dates {
				datesTxt += "   📅 " + d + "\n"
			}
			labelDates := widget.NewLabel(datesTxt)

			headerVille := container.NewBorder(nil, nil, widget.NewLabel("+ "+villePropre), btnMap)

			datesContainer.Add(headerVille)
			datesContainer.Add(labelDates)
			datesContainer.Add(widget.NewSeparator())
		}

		finalContent := container.NewVBox(infoContainer, datesContainer)
		scroll := container.NewVScroll(finalContent)
		scroll.SetMinSize(fyne.NewSize(350, 500))

		var popup *widget.PopUp
		btnClose := widget.NewButton("Fermer", func() {
			popup.Hide()
			list.Unselect(id)
		})

		popupContent := container.NewBorder(nil, btnClose, nil, nil, scroll)
		popup = widget.NewModalPopUp(popupContent, myWindow.Canvas())
		popup.Show()
	}

	filtresBox := container.NewVBox(
		checkFavorites,
		widget.NewSeparator(),
		labelDateCreation, sliderDateMin, sliderDateMax, labelDateValues,
		widget.NewSeparator(),
		labelAlbum, sliderAlbumMin, sliderAlbumMax, labelAlbumValues,
		widget.NewSeparator(),
		labelMembres, sliderMembresMin, sliderMembresMax, labelMembresValues,
		widget.NewSeparator(),
		labelLieu, selectLieu,
		boutonFiltrer,
	)

	filtresScroll := container.NewVScroll(filtresBox)
	filtresScroll.SetMinSize(fyne.NewSize(260, 0))

	contenuPrincipal := container.NewBorder(barreRecherche, nil, nil, nil, list)

	contenuFinal := container.NewBorder(nil, nil, filtresScroll, nil, contenuPrincipal)

	myWindow.SetContent(contenuFinal)
	myWindow.ShowAndRun()
}
