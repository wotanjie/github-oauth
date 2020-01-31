package models

type User struct {
	AvatarUrl string `json:"avatar_url"`
	Name string `json:"name"`
	Email string `json:"email"`
}
