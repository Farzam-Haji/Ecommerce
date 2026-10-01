package repositories

import "errors"

var ErrRowNotFound = errors.New("No Matching Row")
var ErrDuplicateKey = errors.New("Row Already Exist")
var ErrFailUpdate = errors.New("Couldnt Update The Row")

//transaction
var ErrCartNotFound = errors.New("Cart Not Found")
var ErrEmptyCart = errors.New("Cart Is Empty")
var ErrInsufficientInventory = errors.New("Insufficient Inventory")
var ErrProductUnavailable = errors.New("Unavailable Product")

var ErrOrderNotFound = errors.New("Order Not Found")
var ErrNonPending = errors.New("Order Not Pending")