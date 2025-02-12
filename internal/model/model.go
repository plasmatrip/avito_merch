package model

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

// AuthRequest - запрос на аутентификацию
type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse - ответ на запрос на аутентификацию с токеном
type AuthResponse struct {
	Token string `json:"token"`
}

// SendCoinrequest - отправка монеты
type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// ErrorResponse - возвращаемая ошибка
type ErrorResponse struct {
	Errors string `json:"errors"`
}

// InfoResponse - информация о монетах, инвентаре и истории транзакций
type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

// Inventory - информация о товаре
type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinHistory - история получения/передачи монет
type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

// Received - полученная монета
type Received struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

// Sent - отправленная монета
type Sent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// Transaction - транзакция
type Transaction struct {
	Login  string `json:"login"`
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

// Claims - токен+id пользователя
type Claims struct {
	jwt.StandardClaims
	UserdID uuid.UUID
}

// ValidLogin - пустая структура для передачи id пользователя
// в контексте при успешной авторизации
type ValidLogin struct {
}
