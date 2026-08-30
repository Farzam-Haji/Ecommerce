package repositories

import "errors"

var ErrRowNotFound = errors.New("No Matching Row")
var ErrDuplicateKey = errors.New("Row Already Exist")