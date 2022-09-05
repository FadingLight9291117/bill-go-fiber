package model

type BillType int

const (
	BillConsume BillType = iota
	BillIncome
)

type Bill struct {
	ID      string   `json:"id,omitempty" bson:"_id,omitempty"`
	Type    BillType `json:"type"`
	Date    string   `json:"date"`
	Money   float64  `json:"money"`
	Cls     string   `json:"cls"`
	Label   string   `json:"label"`
	Options string   `json:"options,omitempty"`
}
