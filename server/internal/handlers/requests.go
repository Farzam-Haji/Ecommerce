package handlers

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
