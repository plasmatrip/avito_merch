// package db_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/gofrs/uuid"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/plasmatrip/avito_merch/internal/config"
// 	"github.com/plasmatrip/avito_merch/internal/logger"
// 	mock_db "github.com/plasmatrip/avito_merch/internal/storage/db/mock"
// )

// type DBTestSuite struct {
// 	suite.Suite
// 	config *config.Config
// 	logger *logger.Logger
// }

// func (suite *DBTestSuite) SetupSuite() {
// 	var err error
// 	suite.config = &config.Config{
// 		Host:         "localhost:8080",
// 		Database:     "postgres://avito_merch:password@localhost:5432/avito_merch?sslmode=disable",
// 		LogLevel:     "debug",
// 		TokenSecret:  "T0kenS3cRE7",
// 		ReadTimeout:  5,
// 		WriteTimeout: 10,
// 		IdleTimeout:  60,
// 	}
// 	suite.logger, err = logger.NewLogger("debug")
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TestAgentSuite(t *testing.T) {
// 	suite.Run(t, new(DBTestSuite))
// }

// func (suite *DBTestSuite) TestPing() {
// 	ctrl := gomock.NewController(suite.T())
// 	defer ctrl.Finish()

// 	m := mock_db.NewMockDB(ctrl)

// 	suite.NotNil(m)

// 	m.EXPECT().
// 		Ping(gomock.Any()).
// 		Return(nil).Times(1)

// 	suite.NotPanics(func() {
// 		m.Ping(context.Background())
// 	})
// }

// func (suite *DBTestSuite) TestBuyItem() {
// 	ctrl := gomock.NewController(suite.T())
// 	defer ctrl.Finish()

// 	m := mock_db.NewMockDB(ctrl)

// 	suite.NotNil(m)

// 	m.EXPECT().
// 		BuyItem(gomock.Any(), gomock.Any(), gomock.Any()).
// 		Return(nil).Times(1)

// 	err := m.BuyItem(context.Background(), uuid.Nil, "")

// 	suite.Assert().NoError(err)
// }

package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"

	"github.com/plasmatrip/avito_merch/internal/model"
	mock_db "github.com/plasmatrip/avito_merch/internal/storage/db/mock"
	"github.com/stretchr/testify/assert"
)

func TestBuyItemWithMock(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock, err := pgxmock.NewPool()

	// repo := db.PostgresDB{
	// 	DB: mock,
	// }

	// mockDB := mock_db.NewMockDB(ctrl)
	// ctx := context.Background()
	// userID, _ := uuid.NewV4()
	// item := "t-shirt"

	// mockDB.EXPECT().BuyItem(ctx, userID, item).Return(nil)

	// err = mockDB.BuyItem(ctx, userID, item)
	// assert.NoError(t, err)
}

func TestUserAuthWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDB(ctrl)
	ctx := context.Background()
	login := model.AuthRequest{UserName: "user", Password: "pass"}
	mockID, _ := uuid.NewV4()

	mockDB.EXPECT().UserAuth(ctx, login).Return(mockID, nil)

	id, err := mockDB.UserAuth(ctx, login)
	assert.NoError(t, err)
	assert.Equal(t, mockID, id)
}

func TestUserAuthUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDB(ctrl)
	ctx := context.Background()
	login := model.AuthRequest{UserName: "newuser", Password: "newpass"}

	mockDB.EXPECT().UserAuth(ctx, login).Return(uuid.Nil, errors.New("user not found"))

	id, err := mockDB.UserAuth(ctx, login)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestSendCoinWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDB(ctrl)
	ctx := context.Background()
	fromUser, _ := uuid.NewV4()
	sendReq := model.SendCoinRequest{ToUser: "receiver", Amount: 50}

	mockDB.EXPECT().SendCoin(ctx, fromUser, sendReq).Return(nil)

	err := mockDB.SendCoin(ctx, fromUser, sendReq)
	assert.NoError(t, err)
}

func TestInfoWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDB(ctrl)
	ctx := context.Background()
	userID, _ := uuid.NewV4()
	info := model.InfoResponse{Coins: 100}

	mockDB.EXPECT().Info(ctx, userID).Return(info, nil)

	res, err := mockDB.Info(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, 100, res.Coins)
}

func TestPingWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_db.NewMockDB(ctrl)
	ctx := context.Background()

	mockDB.EXPECT().Ping(ctx).Return(nil)

	err := mockDB.Ping(ctx)
	assert.NoError(t, err)
}
