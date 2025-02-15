package db

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"

	"github.com/plasmatrip/avito_merch/internal/storage/db/queries"
)

type DB interface {
	Ping(ctx context.Context) error
	UserAuth(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error)
	BuyItem(ctx context.Context, userID uuid.UUID, item string) error
	SendCoin(ctx context.Context, fromUser uuid.UUID, userSendCoin model.SendCoinRequest) error
	Info(ctx context.Context, userID uuid.UUID) (model.InfoResponse, error)
}

type PostgresDB struct {
	DB  *pgxpool.Pool
	Log logger.Logger
}

func NewRepository(ctx context.Context, dsn string, log logger.Logger) (*PostgresDB, error) {
	// запускаем миграцию
	err := startMigration(dsn)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, err
		} else {
			log.Sugar.Debugw("the database exists, there is nothing to migrate")
		}
	} else {
		log.Sugar.Debugw("database migration was successful")
	}

	// открываем БД
	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		DB:  db,
		Log: log,
	}, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

// StartMigration запускает миграцию
func startMigration(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
		return err
	}
	return nil
}

// Ping проверяет подключение к БД
func (r PostgresDB) Ping(ctx context.Context) error {
	return r.DB.Ping(ctx)
}

// Close закрывает подключение к БД
func (r PostgresDB) Close() {
	r.DB.Close()
}

// func (r PostgresDB) FindUser(ctx context.Context, login model.AuthRequest) (uuid.UUID, error) {
func (r PostgresDB) findUser(ctx context.Context, login model.AuthRequest) (uuid.UUID, error) {
	var user model.AuthRequest
	var userID uuid.UUID

	err := r.DB.QueryRow(ctx, queries.SelectUser, pgx.NamedArgs{"login": login.UserName}).Scan(&userID, &user.UserName, &user.Password)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return userID, apperr.ErrBadLogin
		}
	}

	savedHash, err := hex.DecodeString(user.Password)
	if err != nil {
		return userID, err
	}

	h := sha256.New()
	h.Write([]byte([]byte(login.Password)))
	hash := h.Sum(nil)

	if user.UserName != login.UserName || !bytes.Equal(hash, savedHash) {
		return userID, apperr.ErrBadLogin
	}

	return userID, nil
}

// UserAuth проверка аутентификационных данных, в случае отсутсвия пользователя - регистрация
func (r PostgresDB) UserAuth(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error) {
	userID, err := r.findUser(ctx, userLogin)
	if err == nil {
		return userID, nil
	}

	h := sha256.New()
	h.Write([]byte([]byte(userLogin.Password)))
	hash := hex.EncodeToString(h.Sum(nil))

	var id uuid.UUID

	err = r.DB.QueryRow(ctx, queries.InsertUser, pgx.NamedArgs{
		"date":     time.Now(),
		"login":    userLogin.UserName,
		"password": hash,
	}).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// BuyItem обработка запороса покупки мерча
func (r PostgresDB) BuyItem(ctx context.Context, userID uuid.UUID, item string) error {
	var item_id uuid.UUID
	var item_price int

	err := r.DB.QueryRow(ctx, queries.SelectItem, pgx.NamedArgs{
		"item_name": item,
	}).Scan(&item_id, &item_price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apperr.ErrItemNotFound
		}
		return err
	}

	var user_anount int
	err = r.DB.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&user_anount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apperr.ErrAccountNotFound
		}
		return err
	}

	if user_anount < item_price {
		return apperr.ErrInsufficientFunds
	}

	ct, err := r.DB.Exec(ctx, queries.BuyItem, pgx.NamedArgs{
		"user_id":   userID,
		"item_name": item,
		"price":     item_price,
		"merch_id":  item_id,
	})
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return apperr.ErrMerchNotBought
	}

	return nil
}

// SendCoin обработка запороса отправки монет
func (r PostgresDB) SendCoin(ctx context.Context, fromUser uuid.UUID, sendCoin model.SendCoinRequest) error {
	var user_anount int

	//проверяем наличие счета отправителя
	err := r.DB.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
		"user_id": fromUser,
	}).Scan(&user_anount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apperr.ErrSenderNotFound
		}
		return err
	}

	// проверяем баланс у отправителя
	if user_anount < sendCoin.Amount {
		return apperr.ErrInsufficientFunds
	}

	var toUser uuid.UUID
	//проверяем наличие счета получателя
	err = r.DB.QueryRow(ctx, queries.SelectUserID, pgx.NamedArgs{
		"login": sendCoin.ToUser,
	}).Scan(&toUser)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apperr.ErrRecipientNotFound
		}
		return err
	}

	// проверяем, что отправитель и получатель не один пользователь
	if fromUser == toUser {
		return apperr.ErrSenderAndRecipientAreTheSame
	}

	// начинаем транзакцию
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	// при ошибке коммита откатываем назад
	defer func() error {
		return tx.Rollback(ctx)
	}()

	// обновляем монеты у отправителя
	_, err = tx.Exec(ctx, queries.UpdateCoin, pgx.NamedArgs{
		"user_id": fromUser,
		"amount":  -sendCoin.Amount,
	})
	if err != nil {
		return err
	}

	// обновляем монеты у получателя
	_, err = tx.Exec(ctx, queries.UpdateCoin, pgx.NamedArgs{
		"user_id": toUser,
		"amount":  sendCoin.Amount,
	})
	if err != nil {
		return err
	}

	// записываем информацию о транзакции
	_, err = tx.Exec(ctx, queries.InsertTransaction, pgx.NamedArgs{
		"from_user_id": fromUser,
		"to_user_id":   toUser,
		"amount":       sendCoin.Amount,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Info возвращает информацию о монетах, инвентаре и истории транзакций
func (r PostgresDB) Info(ctx context.Context, userID uuid.UUID) (model.InfoResponse, error) {
	var infoResponse model.InfoResponse

	//проверяем наличие счета пользователя
	err := r.DB.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&infoResponse.Coins)
	if err != nil {
		if err == pgx.ErrNoRows {
			return infoResponse, apperr.ErrAccountNotFound
		}
		return model.InfoResponse{}, err
	}

	//получаем купленный мерч
	rows, err := r.DB.Query(ctx, queries.SelectPurchases, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return model.InfoResponse{}, err
	}

	infoResponse.Inventory = make([]model.Inventory, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		inventory := model.Inventory{}

		err := rows.Scan(
			&inventory.Type,
			&inventory.Quantity,
		)
		if err != nil {
			return model.InfoResponse{}, err
		}

		infoResponse.Inventory = append(infoResponse.Inventory, inventory)
	}
	rows.Close()

	// получаем историю транзакций
	rows, err = r.DB.Query(ctx, queries.SelectTransactions, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return model.InfoResponse{}, err
	}

	for rows.Next() {
		transaction := model.Transaction{}

		err := rows.Scan(
			&transaction.Login,
			&transaction.Amount,
			&transaction.Type,
		)
		if err != nil {
			return model.InfoResponse{}, err
		}

		if transaction.Type == "sent" {
			infoResponse.CoinHistory.Sent = append(infoResponse.CoinHistory.Sent, model.Sent{
				ToUser: transaction.Login,
				Amount: transaction.Amount,
			})
		} else {
			infoResponse.CoinHistory.Received = append(infoResponse.CoinHistory.Received, model.Received{
				FromUser: transaction.Login,
				Amount:   transaction.Amount,
			})
		}
	}

	return infoResponse, nil
}
