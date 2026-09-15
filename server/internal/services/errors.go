package services

import "errors"

//general
var ErrInvalidName = errors.New("invalid name")
var ErrInvalidQuantity = errors.New("invalid quantity")
var ErrInvalidID = errors.New("invalid id")
var ErrInternal= errors.New("internal server issue")
var ErrUnauthorized = errors.New("Action Unauthorized")
var ErrInvalidPrice = errors.New("invalid price")

// product
var ErrProductNotFound = errors.New("product not found")
var ErrDuplicateProduct = errors.New("Product Already Exist")

// user
var ErrEmailUnavailable = errors.New("Email Already Exists")
var ErrWrongCredential = errors.New("one or both credentials are wrong")
var ErrInvalidPassword = errors.New("password must contain at least 8 character")
var ErrTokenGenFailed = errors.New("Failed To Generate Token. Login Again")

//cart
var ErrCartNotFound = errors.New("cart not found")
var ErrInsufficientQuantity = errors.New("product is not available in this quantity")
var ErrItemNotFound = errors.New("item not found")

//order
var ErrOrderNotFound = errors.New("Order Not Found")
var ErrEmptyOrder = errors.New("There Is No Item In the Cart")
