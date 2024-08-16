package models

type User struct {
	UID      int    `json:"u_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
