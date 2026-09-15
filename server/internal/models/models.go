package models

import "time"

type Product struct {
	ID 			int64 	`json:"id"`
	Name 		string 	`json:"name"`
	Quantity	int32 	`json:"quantity"`
	SellerID	int64 	`json:"seller_id"`
	Price		int64 	`json:"price"` 	//minor units
}

type User struct {
	ID 				int64    	`json:"id"`
	Email 			string  	`json:"email"`
	PasswordHash	string  	`json:"password_hash"`
	CreatedAt 		time.Time  	`json:"created_at"`
	UpdatedAt		time.Time  	`json:"updated_at"`
}

type Cart struct {
	ID  		int64  		`json:"id"`
	UserID 		int64  		`json:"user_id"`
	CreatedAt  	time.Time  	`json:"created_at"`
}

type CartItem struct {
	ID 			int64  		`json:"id"`
	CartID		int64  		`json:"cart_id"`
	ProductID	*int64  	`json:"product_id"`
	Name		string  	`json:"name"`
	Quantity	int32  		`json:"quantity"`
	Price  		*int64  	`json:"price"`
}

type Order struct {
	ID  		int64  		`json:"id"`
	UserID  	*int64  	`json:"user_id"`
	Payment 	string 		`json:"payment"`
	CreatedAt   time.Time  	`json:"created_at"`
	UpdatedAt  	time.Time  	`json:"updated_at"`
	ExpiresAt  	time.Time  	`json:"expires_at"`
}

type OrderItem struct {
	ID  		int64  		`json:"id"`
	OrderID   	int64  		`json:"order_id"`
	ProductID  	*int64  		`json:"product_id"`
	Name  		string 		`json:"name"`
	Price  		int64  		`json:"price"`
	Quantity  	int32 		`json:"quantity"`
}