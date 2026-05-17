package handlers


type CreateProductRequest struct {
	Name 	string 		`json:"name" binding:"required"`
	Quantity 	int 	`json:"quantity" binding:"required,gt=0"`
}

type UpdateProductRequest struct {
	// ID 	*int 	`json:"id" binding:"required"`
	Name 	*string 	`json:"name" binding:"required"`
	Quantity *int 	`json:"quantity" binding:"required,gt=0"`
}

type DeleteProductRequest struct {
	ID 			int 		`json:"id" binding:"required"`
	// Name 		string 		`json:"name" binding:"required"`
	// Quantity 	int 		`json:"quantity" binding:"required,gt=0"`
}