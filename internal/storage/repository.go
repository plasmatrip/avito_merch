package storage

import (
	"context"

	"github.com/gofrs/uuid"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"
)

type Repository interface {
	Ping(ctx context.Context) error
	RegisterUser(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error)
	BuyItem(ctx context.Context, userID uuid.UUID, item string) error
}

type PostgresDB struct {
	db  *pgxpool.Pool
	log logger.Logger
}
