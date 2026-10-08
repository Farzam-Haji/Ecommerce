package services

import "errors"

// general
var ErrInvalidID = errors.New("invalid id")
var ErrForbidden = errors.New("authenticated user is not allowed to perform this action")
var ErrInternal = errors.New("internal server error")

// product
var ErrInvalidName = errors.New("invalid product name")
var ErrInvalidQuantity = errors.New("invalid quantity")
var ErrInvalidPrice = errors.New("invalid price")
var ErrProductNotFound = errors.New("product not found")
var ErrDuplicateProduct = errors.New("product already exists")

// user
var ErrEmailUnavailable = errors.New("email is already registered")
var ErrWrongCredential = errors.New("one or both credentials are wrong")
var ErrInvalidPassword = errors.New("password must contain at least 8 characters")
var ErrTokenGenFailed = errors.New("failed to generate authentication token")

// cart
var ErrCartNotFound = errors.New("cart not found")
var ErrInsufficientQuantity = errors.New("product is not available in the requested quantity")
var ErrItemNotFound = errors.New("cart item not found")

// order
var ErrOrderNotFound = errors.New("order not found")
var ErrEmptyOrder = errors.New("cart is empty")
var ErrNonPending = errors.New("order is not pending")

