package models

type Date struct {
	Id 		int 		`json:"id"`
	Dates 	[]string 	`json:"dates"`
}

type DateList struct{
	Indes []Date 		`json:"index"`
}