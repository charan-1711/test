package models

type Todo struct {
	TID    int    `json:"t_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
