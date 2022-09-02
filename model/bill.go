package model

type Bill struct {
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	Date    string `json:"date"`
	Money   string `json:"money"`
	Cls     string `json:"cls"`
	Label   string `json:"label"`
	Options string `json:"options,omitempty"`
}
