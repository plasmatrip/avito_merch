package apperr

import (
	"errors"
	"net/http"
)

var (
	ErrBadLogin                     = errors.New("bad login or password")
	ErrUserNotFound                 = errors.New("user not found")
	ErrLoginAlreadyTaken            = errors.New("login already taken")
	ErrItemNotFound                 = errors.New("item not found")
	ErrInsufficientFunds            = errors.New("insufficient funds")
	ErrAccountNotFound              = errors.New("account not found")
	ErrMerchNotBought               = errors.New("merch not bought")
	ErrSenderNotFound               = errors.New("sender not found")
	ErrRecipientNotFound            = errors.New("recipient not found")
	ErrSenderAndRecipientAreTheSame = errors.New("sender and recipient are the same")

	ErrorMessages = map[error]string{
		ErrItemNotFound:      ErrBadLogin.Error(),
		ErrInsufficientFunds: ErrInsufficientFunds.Error(),
		ErrAccountNotFound:   ErrAccountNotFound.Error(),
		ErrMerchNotBought:    ErrMerchNotBought.Error(),
		ErrRecipientNotFound: ErrRecipientNotFound.Error(),
		ErrSenderNotFound:    ErrSenderNotFound.Error(),
	}

	ErrorStatuses = map[error]int{
		ErrItemNotFound:      http.StatusBadRequest,
		ErrInsufficientFunds: http.StatusBadRequest,
		ErrAccountNotFound:   http.StatusBadRequest,
		ErrMerchNotBought:    http.StatusInternalServerError,
	}
)
