package services

import "errors"

var ErrInvalidName = errors.New("invalid product name")
var ErrInvalidQuantity = errors.New("invalid product quantity")
var ErrInvalidID = errors.New("invalid product id")
var ErrProductNotFound = errors.New("product not found")
var ErrDatabase = errors.New("database issue")
