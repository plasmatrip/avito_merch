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
	err := StartMigration(dsn)
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
func StartMigration(dsn string) error {
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
	return r.db.Ping(ctx)
}

// Close закрывает подключение к БД
func (r PostgresDB) Close() {
	r.db.Close()
}

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

func (r PostgresDB) BuyItem(ctx context.Context, userID uuid.UUID, item string) error {
	var item_id uuid.NullUUID
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

// AddSong добавляет песню
// func (r Repository) AddSong(ctx context.Context, song model.Song) error {
// 	ct, err := r.db.Exec(ctx, queries.AddSong, pgx.NamedArgs{
// 		"group_name":   song.Group,
// 		"song_name":    song.Song,
// 		"release_date": time.Time(song.ReleaseDate),
// 		"lyrics":       song.Text,
// 		"link":         song.Link,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if ct.RowsAffected() == 0 {
// 		r.log.Sugar.Debugw("song not added", "group", song.Group, "song", song.Song)
// 		return errors.New("song not added")
// 	}

// 	return nil
// }

// // DeleteSong удаляет песню
// func (r Repository) DeleteSong(ctx context.Context, song model.Song) error {
// 	ct, err := r.db.Exec(ctx, queries.DeleteSong, pgx.NamedArgs{
// 		"group_name": song.Group,
// 		"song_name":  song.Song,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if ct.RowsAffected() == 0 {
// 		r.log.Sugar.Debugw("song not deleted", "group", song.Group, "song", song.Song)
// 		return errors.New("song not deleted")
// 	}

// 	return nil
// }

// // UpdateSong обновляет песню
// func (r Repository) UpdateSong(ctx context.Context, song model.Song) error {
// 	ct, err := r.db.Exec(ctx, queries.UpdateSong, pgx.NamedArgs{
// 		"group_name":   song.Group,
// 		"song_name":    song.Song,
// 		"release_date": song.ReleaseDate.NilIfZero(),
// 		"lyrics":       song.Text,
// 		"link":         song.Link,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if ct.RowsAffected() == 0 {
// 		r.log.Sugar.Debugw("song not updated", "group", song.Group, "song", song.Song, "release_date", song.ReleaseDate, "lyrics", song.Text, "link", song.Link)
// 		return errors.New("song not updated")
// 	}

// 	return nil
// }

// // GetSongs возвращает список песен по фильтру с пагинацией
// func (r Repository) GetSongs(ctx context.Context, filter *model.Filter) ([]model.Song, error) {
// 	args := []interface{}{}
// 	argID := 1

// 	var query = queries.SelectSongs

// 	if filter.Group != nil {
// 		query += ` AND group_name ILIKE $` + strconv.Itoa(argID)
// 		args = append(args, "%"+*filter.Group+"%")
// 		argID++
// 	}
// 	if filter.Song != nil {
// 		query += ` AND song_name ILIKE $` + strconv.Itoa(argID)
// 		args = append(args, "%"+*filter.Song+"%")
// 		argID++
// 	}
// 	if filter.ReleaseFrom != nil {
// 		query += ` AND release_date >= $` + strconv.Itoa(argID)
// 		args = append(args, *filter.ReleaseFrom)
// 		argID++
// 	}
// 	if filter.ReleaseTo != nil {
// 		query += ` AND release_date <= $` + strconv.Itoa(argID)
// 		args = append(args, *filter.ReleaseTo)
// 		argID++
// 	}
// 	if filter.Text != nil {
// 		query += ` AND lyrics ILIKE $` + strconv.Itoa(argID)
// 		args = append(args, "%"+*filter.Text+"%")
// 		argID++
// 	}
// 	if filter.Link != nil {
// 		query += ` AND link ILIKE $` + strconv.Itoa(argID)
// 		args = append(args, "%"+*filter.Link+"%")
// 		argID++
// 	}

// 	query += ` ORDER BY release_date LIMIT $` + strconv.Itoa(argID)
// 	args = append(args, filter.Limit)
// 	argID++
// 	query += ` OFFSET $` + strconv.Itoa(argID)
// 	args = append(args, filter.Page)

// 	rows, err := r.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var songs []model.Song
// 	for rows.Next() {
// 		var s model.Song
// 		var rd time.Time
// 		err := rows.Scan(&s.Group, &s.Song, &rd, &s.Text, &s.Link)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s.ReleaseDate = model.ReleaseDate(rd)
// 		songs = append(songs, s)
// 	}

// 	return songs, nil
// }

// // GetLyrics возвращает текст песни с пагинацией по куплетам
// // Считаем начало каждого куплета как двойной перевод строки
// func (r Repository) GetLyrics(ctx context.Context, song model.Song, verseNum int) (model.VerseResponse, error) {
// 	var verse model.VerseResponse
// 	var lyrics string

// 	err := r.db.QueryRow(ctx, queries.SelectSong, pgx.NamedArgs{
// 		"group": song.Group,
// 		"song":  song.Song,
// 	}).Scan(&lyrics)

// 	if err != nil {
// 		r.log.Sugar.Debugw("song not found", "group", song.Group, "song", song.Song, "error", err)
// 		return verse, err
// 	}

// 	// Разбиваем на куплеты
// 	verses := strings.Split(lyrics, "\n\n")
// 	totalVerses := len(verses)

// 	if verseNum > totalVerses {
// 		r.log.Sugar.Debugw("verse number out of range", "verse number", verseNum, "group", song.Group, "song", song.Song, "total_verses", totalVerses)
// 		return verse, errors.New("verse number out of range. total verses: " + strconv.Itoa(totalVerses))
// 	}

// 	// Формируем ответ
// 	verse.Group = song.Group
// 	verse.Song = song.Song
// 	verse.Verse = verses[verseNum-1]
// 	verse.VerseNum = verseNum
// 	verse.TotalVerses = totalVerses

// 	return verse, nil
// }
