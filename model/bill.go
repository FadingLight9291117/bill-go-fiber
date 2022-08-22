package model

type Bill struct {
	Date    string `json:"date"`
	Money   string `json:"money"`
	Cls     string `json:"cls"`
	Label   string `json:"label"`
	Options string `json:"options"`
}
