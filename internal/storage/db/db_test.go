package db_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/plasmatrip/avito_merch/internal/storage/db"
	mock_db "github.com/plasmatrip/avito_merch/internal/storage/db/mock"
	"github.com/plasmatrip/avito_merch/internal/storage/db/queries"
)

// var db *sql.DB

type DBTestSuite struct {
	suite.Suite
	mockPool *mock_db.MockPgxPool
	DB       db.PostgresDB
}

func (suite *DBTestSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.mockPool = mock_db.NewMockPgxPool(ctrl)
	suite.DB = db.PostgresDB{
		DB: suite.mockPool,
	}
}

func TestDBSuite(t *testing.T) {
	// run tests
	suite.Run(t, new(DBTestSuite))
}

func (suite *DBTestSuite) TestPing() {
	suite.mockPool.EXPECT().Ping(gomock.Any()).Return(nil)

	suite.NoError(suite.DB.Ping(context.Background()))
}

// func (suite *DBTestSuite) TestBuyItem() {
// 	// commandTag := pgconn.NewCommandTag("")
// 	suite.mockPool.EXPECT().Exec(gomock.Any(), queries.SelectItem, "cup").Return(nil)
// 	suite.mockPool.

// }

func (suite *DBTestSuite) TestBuyItem_Success() {
	ctx := context.Background()
	userID, _ := uuid.NewV4()
	itemID, _ := uuid.NewV4()
	itemPrice, userAmount := 500, 1000

	columns := []string{"id", "price::money::numeric"}
	pgxRow := pgxpoolmo.NewRow(columns, itemID, itemPrice)
	suite.mockPool.EXPECT().QueryRow(ctx, "SELECT id, price::money::numeric FROM merch WHERE name = @item_name", "cup").Return(pgxRow)
	suite.mockPool.EXPECT().QueryRow(ctx, queries.SelectAccount, userID).Return(userAmount)
	suite.mockPool.EXPECT().Exec(ctx, queries.BuyItem, gomock.Any()).Return(pgconn.NewCommandTag("UPDATE 1"), nil)

	suite.NoError(suite.DB.BuyItem(ctx, userID, "item"))
}

// func (suite *DBTestSuite) TestBuyItem_ItemNotFound() {
// 	ctx := context.Background()
// 	userID := uuid.New()
// 	suite.mockPool.EXPECT().QueryRow(ctx, queries.SelectItem, gomock.Any()).Return(nil, pgx.ErrNoRows)
// 	suite.Equal(apperr.ErrItemNotFound, suite.DB.BuyItem(ctx, userID, "item"))
// }
