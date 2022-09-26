package models

type Query struct {
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Order string `json:"order,omitempty"`
}
