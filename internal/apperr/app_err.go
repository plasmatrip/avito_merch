package apperr

import (
	"errors"
	"net/http"
)

var (
	ErrBadLogin                      = errors.New("bad login or password")
	ErrUserNotFound                  = errors.New("user not found")
	ErrLoginAlreadyTaken             = errors.New("login already taken")
	ErrEmptyLoginOrPassword          = errors.New("empty login or password")
	ErrItemNotFound                  = errors.New("item not found")
	ErrInsufficientFunds             = errors.New("insufficient funds")
	ErrAccountNotFound               = errors.New("account not found")
	ErrMerchNotBought                = errors.New("merch not bought")
	ErrSenderNotFound                = errors.New("sender not found")
	ErrRecipientNotFound             = errors.New("recipient not found")
	ErrSenderAndRecipientAreTheSame  = errors.New("sender and recipient are the same")
	ErrAmonutIsLessThanOrEqualToZero = errors.New("amount is less than or equal to zero")
	ErrBadJSON                       = errors.New("bad JSON")
	ErrMecrhNameIsEmpty              = errors.New("merch name is empty")
	ErrInternalServerError           = errors.New("internal server error")
	ErrAuthenticationError           = errors.New("authentication error")
	ErrJWTTotkenError                = errors.New("JWT token error")
	ErrInvalidToken                  = errors.New("invalid token")
	ErrInvalidAuthorizationHeader    = errors.New("invalid authorization header")
	ErrMissingAuthorizationHeader    = errors.New("missing authorization header")

	ErrorMessages = map[error]string{
		ErrItemNotFound:                  ErrBadLogin.Error(),
		ErrInsufficientFunds:             ErrInsufficientFunds.Error(),
		ErrAccountNotFound:               ErrAccountNotFound.Error(),
		ErrMerchNotBought:                ErrMerchNotBought.Error(),
		ErrRecipientNotFound:             ErrRecipientNotFound.Error(),
		ErrSenderNotFound:                ErrSenderNotFound.Error(),
		ErrAmonutIsLessThanOrEqualToZero: ErrAmonutIsLessThanOrEqualToZero.Error(),
		ErrBadJSON:                       ErrBadJSON.Error(),
		ErrSenderAndRecipientAreTheSame:  ErrSenderAndRecipientAreTheSame.Error(),
		ErrMecrhNameIsEmpty:              ErrMecrhNameIsEmpty.Error(),
		ErrEmptyLoginOrPassword:          ErrEmptyLoginOrPassword.Error(),
		ErrAuthenticationError:           ErrAuthenticationError.Error(),
		ErrJWTTotkenError:                ErrJWTTotkenError.Error(),
		ErrInvalidToken:                  ErrInvalidToken.Error(),
		ErrInvalidAuthorizationHeader:    ErrInvalidAuthorizationHeader.Error(),
		ErrMissingAuthorizationHeader:    ErrMissingAuthorizationHeader.Error(),
	}

	ErrorStatuses = map[error]int{
		ErrItemNotFound:                  http.StatusBadRequest,
		ErrInsufficientFunds:             http.StatusBadRequest,
		ErrAccountNotFound:               http.StatusBadRequest,
		ErrMerchNotBought:                http.StatusBadRequest,
		ErrRecipientNotFound:             http.StatusBadRequest,
		ErrSenderNotFound:                http.StatusBadRequest,
		ErrAmonutIsLessThanOrEqualToZero: http.StatusBadRequest,
		ErrMerchNotBought:                http.StatusBadRequest,
		ErrBadJSON:                       http.StatusBadRequest,
		ErrSenderAndRecipientAreTheSame:  http.StatusBadRequest,
		ErrMecrhNameIsEmpty:              http.StatusBadRequest,
		ErrEmptyLoginOrPassword:          http.StatusBadRequest,
		ErrAuthenticationError:           http.StatusConflict,
		ErrJWTTotkenError:                http.StatusUnauthorized,
		ErrInvalidToken:                  http.StatusUnauthorized,
		ErrInvalidAuthorizationHeader:    http.StatusUnauthorized,
		ErrMissingAuthorizationHeader:    http.StatusUnauthorized,
	}
)
