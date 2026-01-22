package main

import (
	"encoding/json"
	"os"
)

var favoriteArtists = make(map[int]bool)
const likesFile = "likes.json"

func LoadLikes() {
	file, err := os.ReadFile(likesFile)
	if err == nil {
		var list []int
		json.Unmarshal(file, &list)
		for _, id := range list {
			favoriteArtists[id] = true
		}
	}
}

func SaveLikes() {
	var list []int
	for id, liked := range favoriteArtists {
		if liked {
			list = append(list, id)
		}
	}
	data, _ := json.Marshal(list)
	os.WriteFile(likesFile, data, 0644)
}

func ToggleLike(artistID int) {
	if favoriteArtists[artistID] {
		delete(favoriteArtists, artistID)
	} else {
		favoriteArtists[artistID] = true
	}
	SaveLikes()
}

func IsLiked(artistID int) bool {
	return favoriteArtists[artistID]
}