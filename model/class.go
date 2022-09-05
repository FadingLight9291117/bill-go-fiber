package model

type Class struct {
	Consume map[string][]string `json:"consume"`
	Income  []string            `json:"income"`
}
