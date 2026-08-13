package models

import "time"

type Product struct {
	ID 		int64 	`json:"id"`
	Name 	string 	`json:"name"`
	Quantity	int32 	`json:"quantity"`
	SellerID	int64 	`json:"seller_id"`
	Price	int64 	`json:"price"` 	//minor units
}

type User struct {
	ID 		int64
	Email 	string
	PasswordHash	string
	CreatedAt 	time.Time
	UpdatedAt	time.Time
}