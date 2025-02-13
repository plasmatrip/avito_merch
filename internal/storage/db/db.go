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

//go:generate mockgen -source=db.go -destination=mock/db.go

type DB interface {
	Ping(ctx context.Context) error
	UserAuth(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error)
	BuyItem(ctx context.Context, userID uuid.UUID, item string) error
	SendCoin(ctx context.Context, fromUser uuid.UUID, userSendCoin model.SendCoinRequest) error
	Info(ctx context.Context, userID uuid.UUID) (model.InfoResponse, error)
}

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
// 	from1, _ := uuid.FromString("c10e162a-3759-4a3d-b36f-4c490e4b3c5a")
// 	from2, _ := uuid.FromString("7bb0fcc6-a318-4898-9563-26e5f55b097c")
// 	from3, _ := uuid.FromString("facdfcdf-cfa0-4e7f-b128-b1e8f1867305")
// 	from4, _ := uuid.FromString("a2cb8bd4-0a6e-41b9-9a8a-dcd5112d3b7b")
// 	from5, _ := uuid.FromString("736080be-2e53-480d-93e1-9f7d78053817")

// 	froms := [5]uuid.UUID{from1, from2, from3, from4, from5}

// 	to1, _ := uuid.FromString("ceb0ebb4-a37f-4a40-ae8b-cdb79f745736")
// 	to2, _ := uuid.FromString("fd1ab2ee-da28-45ff-8d76-82ddc9428e99")
// 	to3, _ := uuid.FromString("17359fd4-32cb-4d3c-b504-a96b1484cbb9")
// 	to4, _ := uuid.FromString("7ba9642b-d519-472c-8674-5e0c3f425492")
// 	to5, _ := uuid.FromString("b7d922ed-75ed-42a2-aaa9-19c6f7866188")
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
// 	uuid1, _ := uuid.FromString("c10e162a-3759-4a3d-b36f-4c490e4b3c5a")
// 	uuid2, _ := uuid.FromString("7bb0fcc6-a318-4898-9563-26e5f55b097c")
// 	uuid3, _ := uuid.FromString("facdfcdf-cfa0-4e7f-b128-b1e8f1867305")
// 	users := [3]uuid.UUID{uuid1, uuid2, uuid3}

// 	item1, _ := uuid.FromString("abca7783-5c28-48ea-a03b-d5bfc5c36404")
// 	item2, _ := uuid.FromString("c575884c-6e01-4018-be12-85f68583d958")
// 	item3, _ := uuid.FromString("4e770e1d-9f8e-4520-a631-31e2d437edca")
// 	item4, _ := uuid.FromString("83f0ceb6-18f7-468b-85b9-6ee26415b0e4")
// 	item5, _ := uuid.FromString("c1c1d9f9-ffb2-4eb5-ac38-6127a971269c")
// 	item6, _ := uuid.FromString("b40b348f-e52f-4278-a18a-d89546017d35")
// 	item7, _ := uuid.FromString("8f79fe40-3beb-4a43-bd75-bdc765342645")
// 	item8, _ := uuid.FromString("42c313b9-0cee-46ce-9981-41ea4f1ed5d0")
// 	item9, _ := uuid.FromString("539d90cb-04e1-4d3e-8a60-bff4ec4a45fb")
// 	item10, _ := uuid.FromString("8d70045a-47c0-4a5a-a87e-388a36b304fb")

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

// func (r PostgresDB) FindUser(ctx context.Context, login model.AuthRequest) (uuid.UUID, error) {
func (r PostgresDB) findUser(ctx context.Context, login model.AuthRequest) (uuid.UUID, error) {
	var user model.AuthRequest
	var userID uuid.UUID

	err := r.db.QueryRow(ctx, queries.SelectUser, pgx.NamedArgs{"login": login.UserName}).Scan(&userID, &user.UserName, &user.Password)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return userID, apperr.ErrBadLogin
		}
	}

	// TODO: восстановление пароля
	// savedHash, err := hex.DecodeString(user.Password)
	// if err != nil {
	// 	return userID, err
	// }

	// h := sha256.New()
	// h.Write([]byte([]byte(login.Password)))
	// hash := h.Sum(nil)

	// if user.UserName != login.UserName || !bytes.Equal(hash, savedHash) {
	// 	return userID, apperr.ErrBadLogin
	// }

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

	err = r.db.QueryRow(ctx, queries.InsertUser, pgx.NamedArgs{
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
	err = r.db.QueryRow(ctx, queries.SelectUserID, pgx.NamedArgs{
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
