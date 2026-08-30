package services

import "errors"

//general
var ErrInvalidName = errors.New("invalid name")
var ErrInvalidQuantity = errors.New("invalid quantity")
var ErrInvalidID = errors.New("invalid id")
var ErrInternal= errors.New("internal server issue")
var ErrInvalidPrice = errors.New("invalid price")

// product
var ErrProductNotFound = errors.New("product not found")
var ErrUnauthorized = errors.New("not allowed to change this product")

// user
var ErrEmailExists = errors.New("email already exists")
var ErrWrongCredential = errors.New("one or both credentials are wrong")
var ErrInvalidPassword = errors.New("password must contain at least 8 character")
var ErrBcrypt = errors.New("bcrypt issue")

//cart
var ErrCartNotFound = errors.New("cart not found")
var ErrInsufficientQuantity = errors.New("product is not available in this quantity")
// var ErrDuplicateItem = errors.New("this item already exist in cart")
var ErrItemNotFound = errors.New("item not found")
