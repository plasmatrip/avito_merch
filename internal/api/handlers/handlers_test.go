package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/plasmatrip/avito_merch/internal/api/handlers"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"
	"github.com/plasmatrip/avito_merch/internal/storage/db"
)

const (
	itemName  = "phone"
	itemPrice = 10
)

type HandlersTestSuite struct {
	suite.Suite
	db          *db.PostgresDB
	pgContainer *postgres.PostgresContainer
	handlers    *handlers.Handlers
}

// настройка окрежения тестов
func (suite *HandlersTestSuite) SetupSuite() {
	ctx := context.Background()

	logger, err := logger.NewLogger(logger.LogLevelDebug)
	if err != nil {
		suite.T().Fatal(err)
	}

	// запускаем контейнер
	pgContainer, err := postgres.Run(ctx,
		"postgres:17.2",
		postgres.WithInitScripts(filepath.Join("..", "..", "storage", "db", "init_test", "init_test_db.sh")),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.pgContainer = pgContainer

	suite.T().Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			suite.T().Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(suite.T(), err, "failed to get connection string")

	stor, err := db.NewRepository(ctx, connStr, *logger)
	assert.NoError(suite.T(), err, "failed to create repository")

	suite.db = stor

	suite.handlers = &handlers.Handlers{
		Stor:   stor,
		Logger: *logger,
	}

	ct, err := suite.db.DB.Exec(ctx, "INSERT INTO merch (id, name, price) VALUES (gen_random_uuid (), @name, @price)", pgx.NamedArgs{
		"name":  itemName,
		"price": itemPrice,
	})
	require.NoError(suite.T(), err, "failed to insert item")
	require.Equal(suite.T(), int64(1), ct.RowsAffected(), "failed to insert item. expected 1 row, got %d", ct.RowsAffected())
}

// удаляем контейнер
func (suite *HandlersTestSuite) TearDownAll(t *testing.T) {
	if err := testcontainers.TerminateContainer(suite.pgContainer); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
	t.Log("TearDownAll")
}

// запускаем тест
func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

// тест на отправку монет
func (suite *HandlersTestSuite) TestSendCoin() {
	ctx := context.Background()

	fromUser := "henry"
	toUser := "john"

	fromUserID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: fromUser,
		Password: "henry",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), fromUserID, "non-zero id was returned")

	toUserID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: toUser,
		Password: "john",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), toUserID, "non-zero id was returned")

	suite.T().Run("coin sent successfully", func(t *testing.T) {
		sc := model.SendCoinRequest{ToUser: toUser, Amount: 100}
		jsonData, _ := jsoniter.Marshal(sc)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(jsonData))
		ctxWV := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: fromUserID})
		req = req.WithContext(ctxWV)
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusOK {
			suite.T().Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	suite.T().Run("sender not found", func(t *testing.T) {
		sc := model.SendCoinRequest{ToUser: toUser, Amount: 100}
		jsonData, _ := jsoniter.Marshal(sc)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(jsonData))
		ctxWV := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: uuid.Nil})
		req = req.WithContext(ctxWV)
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusBadRequest {
			suite.T().Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	suite.T().Run("recipient not found", func(t *testing.T) {
		sc := model.SendCoinRequest{ToUser: "", Amount: 100}
		jsonData, _ := jsoniter.Marshal(sc)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(jsonData))
		ctxWV := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: fromUserID})
		req = req.WithContext(ctxWV)
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusBadRequest {
			suite.T().Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	suite.T().Run("invalid amount", func(t *testing.T) {
		sc := model.SendCoinRequest{ToUser: toUser, Amount: -10}
		jsonData, _ := jsoniter.Marshal(sc)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(jsonData))
		ctxWV := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: uuid.Nil})
		req = req.WithContext(ctxWV)
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusBadRequest {
			suite.T().Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	suite.T().Run("sender and recipient are the same", func(t *testing.T) {
		sc := model.SendCoinRequest{ToUser: toUser, Amount: 10}
		jsonData, _ := jsoniter.Marshal(sc)
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(jsonData))
		ctxWV := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: toUserID})
		req = req.WithContext(ctxWV)
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusBadRequest {
			suite.T().Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	suite.T().Run("invaliud JSON", func(t *testing.T) {
		badData := []byte("{}")
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(badData))
		w := httptest.NewRecorder()

		suite.handlers.SendCoin(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// тест на покупку
func (suite *HandlersTestSuite) TestBuy() {
	ctx := context.Background()

	fromUser := "henry"
	toUser := "john"

	fromUserID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: fromUser,
		Password: "henry",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), fromUserID, "non-zero id was returned")

	toUserID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: toUser,
		Password: "john",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), toUserID, "non-zero id was returned")

	suite.T().Run("buy item successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/"+itemName, nil)
		ctx := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: fromUserID})
		req = req.WithContext(ctx)
		req.SetPathValue("item", itemName)
		w := httptest.NewRecorder()

		suite.handlers.Buy(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	suite.T().Run("buy item fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/buy/hummer", nil)
		ctx := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: fromUserID})
		req = req.WithContext(ctx)
		req.SetPathValue("item", "hummer")
		w := httptest.NewRecorder()

		suite.handlers.Buy(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func (suite *HandlersTestSuite) TestInfo() {
	ctx := context.Background()

	user1 := "boris"

	user2 := "kevin"
	user2Amount := 10

	user3 := "ivan"
	user3Amount := 25

	user1InitialAmount := 1000

	user1ID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: user1,
		Password: user1,
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), user1ID, "non-zero id was returned")

	user2ID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: user2,
		Password: user2,
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), user2ID, "non-zero id was returned")

	user3ID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: user3,
		Password: user3,
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), user3ID, "non-zero id was returned")

	err = suite.db.BuyItem(ctx, user1ID, itemName)
	assert.NoError(suite.T(), err, "an error occurred when buying an existing item")

	err = suite.db.SendCoin(ctx, user1ID, model.SendCoinRequest{
		ToUser: user3,
		Amount: user3Amount,
	})
	assert.NoError(suite.T(), err, "an error occurred while sending coins between existing users")

	err = suite.db.SendCoin(ctx, user2ID, model.SendCoinRequest{
		ToUser: user1,
		Amount: user2Amount,
	})
	assert.NoError(suite.T(), err, "an error occurred while sending coins between existing users")

	suite.T().Run("info received successfully", func(t *testing.T) {
		expected := model.InfoResponse{
			Coins: user1InitialAmount - itemPrice + user2Amount - user3Amount,
			Inventory: []model.Inventory{
				{
					Type:     itemName,
					Quantity: 1,
				},
			},
			CoinHistory: model.CoinHistory{
				Received: []model.Received{
					{
						FromUser: user2,
						Amount:   user2Amount,
					},
				},
				Sent: []model.Sent{
					{
						ToUser: user3,
						Amount: user3Amount,
					},
				},
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		ctx := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: user1ID})
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		suite.handlers.Info(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var got model.InfoResponse
		if err := jsoniter.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		assert.Equal(t, expected, got, "info response does not match expected")
	})

	suite.T().Run("user not found", func(t *testing.T) {
		uuid, err := uuid.NewV4()
		assert.NoError(suite.T(), err, "an error occurred while generating a UUID")

		expected := model.InfoResponse{}

		req := httptest.NewRequest(http.MethodGet, "/info", nil)
		ctx := context.WithValue(req.Context(), model.ValidLogin{}, &model.Claims{UserdID: uuid})
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		suite.handlers.Info(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var got model.InfoResponse
		if err := jsoniter.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Errorf("failed to decode response: %v", err)
		}

		assert.Equal(t, expected, got, "info response does not match expected")
	})
}
