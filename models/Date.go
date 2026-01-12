package models

type Date struct {
	Id 		int 		`json:"id"`
	Dates 	[]string 	`json:"dates"`
}

type Datelist struct{
	Indes []Date 		`json:"index"`
}