package models

import "time"

type Product struct {
	ID 		int 	`json:"id"`
	Name 	string 	`json:"name"`
	Quantity	int 	`json:"quantity"`
}

type User struct {
	ID 		int64
	Email 	string
	PasswordHash	string
	CreatedAt 	time.Time
	UpdatedAt	time.Time
}