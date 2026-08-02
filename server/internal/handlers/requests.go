package handlers

//Product
type CreateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type PatchProductRequest struct {
	Name     *string `json:"name"`
	Quantity *int    `json:"quantity"`
}

// User
type RegisterRequest struct {
	Email string 	`json:"email" binding:"required"`
	Password 	string 	`json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email 	string 	`json:"email" binding:"required"`
	Password	string 	`json:"password" binding:"required"`
}