package db

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
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

//go:generate mockgen -source=db.go -destination=mock/dbmock.go

type DB interface {
	Ping(ctx context.Context) error
	UserAuth(ctx context.Context, userLogin model.AuthRequest) (uuid.UUID, error)
	BuyItem(ctx context.Context, userID uuid.UUID, item string) error
	SendCoin(ctx context.Context, fromUser uuid.UUID, userSendCoin model.SendCoinRequest) error
	Info(ctx context.Context, userID uuid.UUID) (model.InfoResponse, error)
}

type PostgresDB struct {
	DB PgxPool
	// DB  *pgxpool.Pool
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

	// go generateUsers(db)
	// go generateTransactions(db)
	// go geratePurchases(db)

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

func generateTransactions(db *pgxpool.Pool) {
	from1, _ := uuid.FromString("114d785b-ad37-4288-a21a-731a61681c0c")
	from2, _ := uuid.FromString("d3794488-b427-4f60-bc85-7d58e34ab851")
	from3, _ := uuid.FromString("55e5dd08-f782-43da-8ed4-e96712245047")
	from4, _ := uuid.FromString("2cad6224-011d-42ef-acc1-4818680e2986")
	from5, _ := uuid.FromString("fdbb4088-3c63-4cef-a6eb-91ace3e935d4")

	froms := [5]uuid.UUID{from1, from2, from3, from4, from5}

	to1, _ := uuid.FromString("bf6a36b8-ed41-4f06-a36d-d3d73049a9b2")
	to2, _ := uuid.FromString("9eac4a66-91de-4d6d-ac30-b24e520abc39")
	to3, _ := uuid.FromString("e2d2ac49-63d8-4f74-82a2-158d6f8f824f")
	to4, _ := uuid.FromString("b30c374f-8c76-42b7-9e5a-cf362e29f0b3")
	to5, _ := uuid.FromString("73b2822a-63c7-4fee-a496-e9e9a0a35070")

	tos := [5]uuid.UUID{to1, to2, to3, to4, to5}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 4; i < 100000; i++ {
		go func() {
			_, _ = db.Exec(context.Background(), queries.InsertTransaction, pgx.NamedArgs{
				"from_user_id": froms[r.Intn(5)],
				"to_user_id":   tos[r.Intn(5)],
				"amount":       100,
			})
		}()
	}
}

func generateUsers(db *pgxpool.Pool) {
	for i := 4; i < 100000; i++ {
		go func() {
			_, err := db.Exec(context.Background(), queries.InsertUser, pgx.NamedArgs{
				"date":     time.Now(),
				"login":    fmt.Sprintf("user%d", i),
				"password": "password",
			})
			if err != nil {
				fmt.Println(err)
			}
		}()

	}
}

func geratePurchases(db *pgxpool.Pool) {
	uuid1, _ := uuid.FromString("bf6a36b8-ed41-4f06-a36d-d3d73049a9b2")
	uuid2, _ := uuid.FromString("55e5dd08-f782-43da-8ed4-e96712245047")
	uuid3, _ := uuid.FromString("2cad6224-011d-42ef-acc1-4818680e2986")
	users := [3]uuid.UUID{uuid1, uuid2, uuid3}

	item1, _ := uuid.FromString("bfb785ce-8cac-4ab2-8ddf-084506b0c8ce")
	item2, _ := uuid.FromString("224eaee5-9256-4da1-9c20-90bd4732f244")
	item3, _ := uuid.FromString("35929711-366a-477d-8018-4dd6fe986aa4")
	item4, _ := uuid.FromString("4c6e4f99-6522-4081-bff0-3e4cb60a76d3")
	item5, _ := uuid.FromString("445efaa9-70e8-47e8-8f7f-4ab51e5fd5f8")
	item6, _ := uuid.FromString("b01f192e-27d5-4da8-9a7d-95c02e05b5f3")
	item7, _ := uuid.FromString("23e31eb4-c5b6-4ce0-9a59-9b73a6f85eb9")
	item8, _ := uuid.FromString("015a03cc-ec64-4273-9732-46b5dcbfa9fb")
	item9, _ := uuid.FromString("127c5d13-c488-4558-9e16-7d1393cf239b")
	item10, _ := uuid.FromString("dfaec662-c1b2-467d-b97c-0431792a59a7")

	items := [10]uuid.UUID{item1, item2, item3, item4, item5, item6, item7, item8, item9, item10}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100000; i++ {
		go func() {
			db.Exec(context.Background(), `INSERT INTO purchases (id, date, user_id, merch_id)
			values (
				gen_random_uuid (),
				CURRENT_TIMESTAMP,
				@user_id,
				@merch_id
			)`, pgx.NamedArgs{
				"user_id":  users[r.Intn(3)],
				"merch_id": items[r.Intn(10)],
			})
		}()
	}
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

	//проверяем наличие счета отправителя
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
