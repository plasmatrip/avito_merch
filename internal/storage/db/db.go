package db

import (
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

type PostgresDB struct {
	db  *pgxpool.Pool
	log logger.Logger
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
		db:  db,
		log: log,
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

// func generateTransactions(db *pgxpool.Pool) {
// 	from1, _ := uuid.FromString("a0b8bd72-a2c6-4d37-a841-349feee0a8ba")
// 	from2, _ := uuid.FromString("ce39d177-7314-4efe-b680-3b6997b8f49f")
// 	from3, _ := uuid.FromString("5bb43b7b-e75f-45cb-99c8-537c49bde3d5")
// 	from4, _ := uuid.FromString("defccd78-ffcc-4215-900f-a72c17c072e5")
// 	from5, _ := uuid.FromString("50cfb64c-0951-4daa-b39a-b7db51684acc")
// 	froms := [5]uuid.UUID{from1, from2, from3, from4, from5}

// 	to1, _ := uuid.FromString("d13178bb-5e41-4518-a576-a9fa888fcbd3")
// 	to2, _ := uuid.FromString("2dd61419-d1d8-4d5e-83ac-4bf4046b21a8")
// 	to3, _ := uuid.FromString("652e2962-4264-4460-8bfa-cd18af29c100")
// 	to4, _ := uuid.FromString("ff3dd987-ad2e-4402-8d7a-1987c6dea5d4")
// 	to5, _ := uuid.FromString("bb3c6ef1-415b-4ca2-9528-87c58f1f1b2b")
// 	tos := [5]uuid.UUID{to1, to2, to3, to4, to5}

// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	for i := 4; i < 100000; i++ {
// 		go func() {
// 			_, _ = db.Exec(context.Background(), queries.InsertTransaction, pgx.NamedArgs{
// 				"from_user_id": froms[r.Intn(5)],
// 				"to_user_id":   tos[r.Intn(5)],
// 				"amount":       100,
// 			})
// 		}()
// 	}
// }

// func generateUsers(db *pgxpool.Pool) {
// 	for i := 4; i < 100000; i++ {
// 		go func() {
// 			_, err := db.Exec(context.Background(), queries.InsertUser, pgx.NamedArgs{
// 				"date":     time.Now(),
// 				"login":    fmt.Sprintf("user%d", i),
// 				"password": "password",
// 			})
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}()

// 	}
// }

// func geratePurchases(db *pgxpool.Pool) {
// 	uuid1, _ := uuid.FromString("a0b8bd72-a2c6-4d37-a841-349feee0a8ba")
// 	uuid2, _ := uuid.FromString("5bb43b7b-e75f-45cb-99c8-537c49bde3d5")
// 	uuid3, _ := uuid.FromString("ce39d177-7314-4efe-b680-3b6997b8f49f")
// 	users := [3]uuid.UUID{uuid1, uuid2, uuid3}

// 	item1, _ := uuid.FromString("d4356121-78a1-42b3-bf5b-d48a7f7173b3")
// 	item2, _ := uuid.FromString("3399ed21-02d8-4c06-859d-e220f66e5654")
// 	item3, _ := uuid.FromString("13ade107-8e03-4dbb-9bd3-1c61fe8ea77f")
// 	item4, _ := uuid.FromString("47cfb7c5-3372-4ef7-a551-4162b3db60ac")
// 	item5, _ := uuid.FromString("d997a3ab-8164-4b33-b533-266f3e65ca18")
// 	item6, _ := uuid.FromString("b1aa6642-739a-4af1-8f42-fe7415841989")
// 	item7, _ := uuid.FromString("74194691-fbde-4a60-b79c-c87626bc5de3")
// 	item8, _ := uuid.FromString("f9a2c0cd-dca1-4965-a770-43d9c5253deb")
// 	item9, _ := uuid.FromString("364bac57-2f6f-4d22-9c7f-06fd710f9ca6")
// 	item10, _ := uuid.FromString("c153167f-7b84-4672-95ae-4ede74f56b46")

// 	items := [10]uuid.UUID{item1, item2, item3, item4, item5, item6, item7, item8, item9, item10}

// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	for i := 0; i < 100000; i++ {
// 		go func() {
// 			db.Exec(context.Background(), `INSERT INTO purchases (id, date, user_id, merch_id)
// 			values (
// 				gen_random_uuid (),
// 				CURRENT_TIMESTAMP,
// 				@user_id,
// 				@merch_id
// 			)`, pgx.NamedArgs{
// 				"user_id":  users[r.Intn(3)],
// 				"merch_id": items[r.Intn(10)],
// 			})
// 		}()
// 	}
// }

// Ping проверяет подключение к БД
func (r PostgresDB) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

// Close закрывает подключение к БД
func (r PostgresDB) Close() {
	r.db.Close()
}

// RegisterUser регистрация пользователя
func (r PostgresDB) RegisterUser(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error) {
	h := sha256.New()
	h.Write([]byte([]byte(userLogin.Password)))
	hash := hex.EncodeToString(h.Sum(nil))

	var id uuid.UUID

	err := r.db.QueryRow(ctx, queries.InsertUser, pgx.NamedArgs{
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

	err := r.db.QueryRow(ctx, queries.SelectItem, pgx.NamedArgs{
		"item_name": item,
	}).Scan(&item_id, &item_price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return apperr.ErrItemNotFound
		}
		return err
	}

	var user_anount int
	err = r.db.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
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

	ct, err := r.db.Exec(ctx, queries.BuyItem, pgx.NamedArgs{
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
	err := r.db.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
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
	err = r.db.QueryRow(ctx, queries.SelectUser, pgx.NamedArgs{
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
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
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

	//проверяем наличие счета отправителя
	err := r.db.QueryRow(ctx, queries.SelectAccount, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&infoResponse.Coins)
	if err != nil {
		if err == pgx.ErrNoRows {
			return infoResponse, apperr.ErrSenderNotFound
		}
		return model.InfoResponse{}, err
	}

	//получаем купленный мерч
	rows, err := r.db.Query(ctx, queries.SelectPurchases, pgx.NamedArgs{
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
	rows, err = r.db.Query(ctx, queries.SelectTransactions, pgx.NamedArgs{
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
