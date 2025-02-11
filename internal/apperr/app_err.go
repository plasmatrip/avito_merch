package apperr

import (
	"errors"
)

var (
	ErrBadLogin                     = errors.New("bad login or password")
	ErrLoginAlreadyTaken            = errors.New("login already taken")
	ErrItemNotFound                 = errors.New("item not found")
	ErrInsufficientFunds            = errors.New("insufficient funds")
	ErrAccountNotFound              = errors.New("account not found")
	ErrMerchNotBought               = errors.New("merch not bought")
	ErrSenderNotFound               = errors.New("sender not found")
	ErrRecipientNotFound            = errors.New("recipient not found")
	ErrSenderAndRecipientAreTheSame = errors.New("sender and recipient are the same")
)
