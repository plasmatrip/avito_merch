package apperr

import "errors"

var (
	ErrBadLogin          = errors.New("bad login or password")
	ErrLoginAlreadyTaken = errors.New("login already taken")

	ErrItemNotFound      = errors.New("item not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrMerchNotBought    = errors.New("merch not bought")
)
