package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/plasmatrip/avito_merch/internal/api/handlers"
	"github.com/plasmatrip/avito_merch/internal/config"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/router"
	"github.com/plasmatrip/avito_merch/internal/storage/db"
)

func main() {
	// для грейсфул шатдауна слушаем сигнал ОС
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// инициализируем логгер
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Close()

	// инициализируем подключение к БД
	db, err := db.NewRepository(ctx, cfg.Database, *log)
	if err != nil {
		log.Sugar.Infow("database connection error: ", err)
		os.Exit(1)
	}
	defer db.Close()

	//пингуем базу
	if err := db.Ping(ctx); err != nil {
		log.Sugar.Infow("database connection error: ", err)
		os.Exit(1)
	}

	// запускаем веб-сервер
	server := http.Server{
		Addr:         cfg.Host,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler: func(next http.Handler) http.Handler {
			log.Sugar.Infow("The Avito merch store is running. ", "Server address", cfg.Host)
			return next
		}(router.NewRouter(*cfg, *log, *db, handlers.NewHandlers(*cfg, *log, *db))),
	}

	go server.ListenAndServe()

	// ждем сигнал ОС
	<-ctx.Done()

	server.Shutdown(context.Background())

	log.Sugar.Infow("The server has been shut down gracefully")

	os.Exit(0)
}
