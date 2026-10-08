package handlers

import _"github.com/gin-gonic/gin/binding"

//Product
//Price intented to be used as minor units
type CreateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity *int32    `json:"quantity" binding:"required,gte=0"`
	Price 	*int64 	`json:"price" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity *int32    `json:"quantity" binding:"required,gte=0"`
	Price 	*int64 	`json:"price" binding:"required,gte=0"`
}

type PatchProductRequest struct {
	Name     *string `json:"name"`
	Quantity *int32    `json:"quantity" binding:"omitempty,gte=0"`
	Price  	*int64 	`json:"price" binding:"omitempty,gte=0"`
}

// User
type RegisterRequest struct {
	Email string 	`json:"email" binding:"required,email"`
	Password 	string 	`json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email 	string 	`json:"email" binding:"required,email"`
	Password	string 	`json:"password" binding:"required"`
}

// Cart
type AddItemRequest struct {
	ProductID 	int64  		`json:"product_id" binding:"required,gt=0"`
	Quantity  	*int32  		`json:"quantity" binding:"omitempty,gte=0"`
}

type UpdateItemQuantityRequest struct {
	Quantity  	*int32  		`json:"quantity" binding:"omitempty,gte=0"`
}

// Order
type PaymentResultRequest struct {
	PaymentID 		int64 		`json:"payment_id" binding:"required,gt=0"`
	OrderID 		int64 		`json:"order_id" binding:"required,gt=0"`
	Status 			string 		`json:"status"`
}