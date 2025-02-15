package db_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/plasmatrip/avito_merch/internal/apperr"
	"github.com/plasmatrip/avito_merch/internal/logger"
	"github.com/plasmatrip/avito_merch/internal/model"
	"github.com/plasmatrip/avito_merch/internal/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	itemName  = "phone"
	itemPrice = 10
)

type DBTestSuite struct {
	suite.Suite
	db          *db.PostgresDB
	pgContainer *postgres.PostgresContainer
}

func (suite *DBTestSuite) SetupSuite() {
	ctx := context.Background()

	logger, err := logger.NewLogger(logger.LogLevelDebug)
	if err != nil {
		suite.T().Fatal(err)
	}

	pgContainer, err := postgres.Run(ctx,
		"postgres:17.2",
		postgres.WithInitScripts(filepath.Join("init_test", "init_test_db.sh")),
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

	ct, err := suite.db.DB.Exec(ctx, "INSERT INTO merch (id, name, price) VALUES (gen_random_uuid (), @name, @price)", pgx.NamedArgs{
		"name":  itemName,
		"price": itemPrice,
	})
	require.NoError(suite.T(), err, "failed to insert item")
	require.Equal(suite.T(), int64(1), ct.RowsAffected(), "failed to insert item. expected 1 row, got %d", ct.RowsAffected())
}

func (suite *DBTestSuite) TearDownAll(t *testing.T) {
	if err := testcontainers.TerminateContainer(suite.pgContainer); err != nil {
		t.Fatalf("failed to terminate container: %s", err)
	}
	t.Log("TearDownAll")
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) TestPing() {
	err := suite.db.Ping(context.Background())
	assert.NoError(suite.T(), err, "database connection check failed")
}

func (suite *DBTestSuite) TestUserAuth() {
	ctx := context.Background()

	uuid, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: "henry",
		Password: "henry",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), uuid, "non-zero id was returned")

}

func (suite *DBTestSuite) TestBuyItemSuccess() {
	ctx := context.Background()

	userID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: "henry",
		Password: "henry",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), userID, "non-zero id was returned")

	err = suite.db.BuyItem(ctx, userID, itemName)
	assert.NoError(suite.T(), err, "an error occurred when buying an existing item")
}

func (suite *DBTestSuite) TestBuyItemFail() {
	ctx := context.Background()

	userID, err := suite.db.UserAuth(ctx, model.AuthRequest{
		UserName: "henry",
		Password: "henry",
	})
	assert.NoError(suite.T(), err, "an error occurred during user authorization")
	assert.NotNil(suite.T(), userID, "non-zero id was returned")

	err = suite.db.BuyItem(ctx, userID, "hummer")
	assert.Error(suite.T(), err, "an error occurred when trying to buy a non-existent item")
	assert.ErrorIs(suite.T(), err, apperr.ErrItemNotFound, "an unexpected error occurred when attempting to purchase a non-existent item")
}

func (suite *DBTestSuite) TestSendCoin() {
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
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), toUserID)

	testCases := []struct {
		testName   string
		fromUserID uuid.UUID
		toUser     string
		amount     int
		wantErr    error
		f          func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool
	}{
		{
			testName:   "coin sent successfully",
			fromUserID: fromUserID,
			toUser:     toUser,
			amount:     10,
			f: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.NoError(t, err, "an error occurred while sending coins between existing users")
			},
		},
		{
			testName:   "sender not found",
			fromUserID: uuid.Nil,
			toUser:     toUser,
			amount:     10,
			f: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, apperr.ErrSenderNotFound, "when sending coins from a non-existent user, an unexpected error occurred")
			},
		},
		{
			testName:   "recipient not found",
			fromUserID: fromUserID,
			toUser:     "",
			amount:     10,
			f: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, apperr.ErrRecipientNotFound, "when sending coins to a non-existent user, an unexpected error occurred")
			},
		},
		{
			testName:   "insufficient funds",
			fromUserID: fromUserID,
			toUser:     toUser,
			amount:     1000000,
			f: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, apperr.ErrInsufficientFunds, "when sending more coins than the user has, an unexpected error occurred")
			},
		},
		{
			testName:   "sender and recipient are the same",
			fromUserID: fromUserID,
			toUser:     fromUser,
			amount:     10,
			f: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.ErrorIs(t, err, apperr.ErrSenderAndRecipientAreTheSame, "unexpected error when a user sends coins to themselves")
			},
		},
	}

	for _, test := range testCases {
		suite.T().Run(test.testName, func(t *testing.T) {
			err = suite.db.SendCoin(ctx, test.fromUserID, model.SendCoinRequest{
				ToUser: test.toUser,
				Amount: test.amount,
			})
			test.f(t, err)
		})
	}
}

func (suite *DBTestSuite) TestInfo() {
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

	info, err := suite.db.Info(ctx, user1ID)
	assert.NoError(suite.T(), err, "an error occurred while getting information about a user")
	assert.Equal(suite.T(), info.Coins, user1InitialAmount-itemPrice+user2Amount-user3Amount, "unexpected number of coins")
	assert.Equal(suite.T(), info.Inventory[0].Type, itemName, "unexpected type of items")
	assert.Equal(suite.T(), info.Inventory[0].Quantity, 1, "unexpected quantity of items")
	assert.Equal(suite.T(), info.CoinHistory.Received[0].FromUser, user2, "unexpected sender of coins. %+v", info.CoinHistory)
	assert.Equal(suite.T(), info.CoinHistory.Received[0].Amount, user2Amount, "unexpected amount of received coins")
	assert.Equal(suite.T(), info.CoinHistory.Sent[0].ToUser, user3, "unexpected recipient of coins")
	assert.Equal(suite.T(), info.CoinHistory.Sent[0].Amount, user3Amount, "unexpected amount of sent coins")

	uuid, err := uuid.NewV4()
	assert.NoError(suite.T(), err, "an error occurred while generating a UUID")

	emptyInfo := model.InfoResponse{}
	info, err = suite.db.Info(ctx, uuid)
	assert.ErrorIs(suite.T(), err, apperr.ErrAccountNotFound, "unexpected error when getting information about a non-existent user")
	assert.Equal(suite.T(), emptyInfo, info, "unexpected response when getting information about a non-existent user")
}
