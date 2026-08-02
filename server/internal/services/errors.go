package services

import "errors"

// product
var ErrInvalidName = errors.New("invalid product name")
var ErrInvalidQuantity = errors.New("invalid product quantity")
var ErrInvalidID = errors.New("invalid product id")
var ErrProductNotFound = errors.New("product not found")
var ErrInternal= errors.New("internal server issue")

// user
var ErrEmailExists = errors.New("email already exists")
var ErrWrongCredential = errors.New("one or both credentials are wrong")
var ErrInvalidPassword = errors.New("password must contain at least 8 character")
var ErrBcrypt = errors.New("bcrypt issue")