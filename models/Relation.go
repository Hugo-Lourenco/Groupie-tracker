package models

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationList struct {
	Index []Relation `json:"index"`
}