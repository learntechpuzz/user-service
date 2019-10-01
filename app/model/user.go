package model

// User represents user domain model
type User struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
}
