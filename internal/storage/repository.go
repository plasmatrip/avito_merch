package storage

import (
	"context"

	"github.com/gofrs/uuid"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/plasmatrip/avito_merch/internal/model"
)

type Repository interface {
	Ping(ctx context.Context) error
	UserAuth(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error)
	BuyItem(ctx context.Context, userID uuid.UUID, item string) error
	SendCoin(ctx context.Context, fromUser uuid.UUID, userSendCoin model.SendCoinRequest) error
	Info(ctx context.Context, userID uuid.UUID) (model.InfoResponse, error)
}
